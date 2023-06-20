package slogt_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/slog"

	"github.com/pyd/slogt"
)

/*
Testing the slogt.Log struct.
*/

var _ = Describe("The Log object", func() {

	var observer *slogt.Observer
	var logger *slog.Logger
	var log slogt.Log

	BeforeEach(func() {
		observer = new(slogt.Observer)
		handler, _ := slogt.NewDefaultObserverHandler(observer)
		logger = slog.New(handler)
	})

	Describe("provides message, time and level getters.", func() {

		var logMessage string
		var logLevel slog.Level

		BeforeEach(func() {
			logMessage = "log message"
			logLevel = slog.LevelWarn
			logger.Log(context.Background(), logLevel, logMessage)
			var logFound bool
			log, logFound = observer.FindLog(1)
			Expect(logFound).To(BeTrue())
		})

		It("should return the message of the log", func() {
			Expect(log.Message()).To(Equal(logMessage))
		})

		It("should return the time when the log was created", func() {
			Expect(log.Time()).To(BeTemporally("~", time.Now(), time.Millisecond*500))
		})

		It("should return the level of the log", func() {
			Expect(log.Level()).To(Equal(logLevel))
		})
	})

	Describe("provides an attribute finder by key.", func() {

		var attribute slog.Attr
		var attributeFound bool
		var searchedKey string

		BeforeEach(func() {
			logger.Error(
				"error message",
				slog.String("username", "gopher"),
				slog.Group("request", slog.String("method", "POST"), slog.Bool("secured", true)),
				slog.Group("user", slog.Group("profile", slog.Int("age", 22), slog.Bool("admin", true))),
			)
			var logFound bool
			log, logFound = observer.FindLog(1)
			Expect(logFound).To(BeTrue())
		})

		JustBeforeEach(func() {
			attribute, attributeFound = log.FindAttribute(searchedKey)
		})

		Context("If attribute is root (not nested)", func() {

			When("there is a match", func() {

				BeforeEach(func() {
					searchedKey = "username"
				})

				It("should return attribute found", func() {
					Expect(attributeFound).To(BeTrue())
				})

				It("should return the matching attribute", func() {
					Expect(attribute.Value.String()).To(Equal("gopher"))
				})
			})

			When("there is no match", func() {

				BeforeEach(func() {
					searchedKey = "unknownKey"
				})

				It("should return attribute not found", func() {
					Expect(attributeFound).To(BeFalse())
				})

				It("should return a zero-ed attribute", func() {
					Expect(attribute.Value.Any()).To(BeNil())
				})
			})
		})

		Context("If attribute is nested", func() {

			When("there is a match", func() {

				BeforeEach(func() {
					searchedKey = "user.profile.age"
				})

				It("should return attribute found", func() {
					Expect(attributeFound).To(BeTrue())
				})

				It("should return the matching attribute", func() {
					Expect(attribute.Value.Int64()).To(Equal(int64(22)))
				})
			})

			When("there is no match", func() {

				BeforeEach(func() {
					searchedKey = "user.profile.unknownkey"
				})

				It("should return attribute not found", func() {
					Expect(attributeFound).To(BeFalse())
				})

				It("should return a zero-ed attribute", func() {
					Expect(attribute.Value.Any()).To(BeNil())
				})
			})
		})

	})

})
