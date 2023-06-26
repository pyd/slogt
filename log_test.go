package slogt_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/slog"

	"github.com/pyd/slogt"
)

/*
Testing the slogt.Log struct.

Note: GetBuilInAttributes() and  GetSharedAttributes() are implicitely tested by the finders tests.
*/
var _ = Describe("The Log struct", func() {

	var observer *slogt.Observer
	var handler slogt.ObserverHandler
	var logger *slog.Logger
	var log slogt.Log

	BeforeEach(func() {
		observer = new(slogt.Observer)
		handler, _ = slogt.NewDefaultObserverHandler(observer)
		logger = slog.New(handler)
	})

	Describe("has getters for message, time and level.", func() {

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

	Describe("has a shared attribute finder", func() {

		var attribute slog.Attr
		var attributeFound bool

		BeforeEach(func() {
			logger = logger.With(slog.Int("userId", 47))
			logger.Error("error message")
			var logFound bool
			log, logFound = observer.FindLog(1)
			Expect(logFound).To(BeTrue())
		})

		When("attribute exists", func() {

			BeforeEach(func() {
				attribute, attributeFound = log.FindSharedAttribute("userId")
			})

			It("should return attribute found", func() {
				Expect(attributeFound).To(BeTrue())
			})

			It("should return matching attribute", func() {
				Expect(attribute.Value.Int64()).To(Equal(int64(47)))
			})
		})

		When("attribute does not exist", func() {

			BeforeEach(func() {
				attribute, attributeFound = log.FindSharedAttribute("unknownKey")
			})

			It("should return attribute not found", func() {
				Expect(attributeFound).To(BeFalse())
			})

			It("should return a zero-ed attribute", func() {
				Expect(attribute.Value.Any()).To(BeNil())
			})
		})
	})

	DescribeTable("has a builtin attribute finder.",
		func(key string, attributeFound bool, attributeValue slog.Value) {

			logger.Error(
				"error message",
				slog.String("client", "frontend"),
				slog.Group("request", slog.String("method", "POST"), slog.Bool("secured", true)),
			)

			log, logFound := observer.FindLog(1)
			Expect(logFound).To(BeTrue())
			attr, attrFound := log.FindBuiltInAttribute(key)

			Expect(attrFound).To(Equal(attributeFound), fmt.Sprintf("attribute found should be %t", attributeFound))
			Expect(attr.Value).To(Equal(attributeValue), fmt.Sprintf("attribute value should be %v", attributeValue))
		},
		Entry("there is an attribute matching this single key",
			"client", true, slog.StringValue("frontend")),

		Entry("there is no attribute matching this single key",
			"unknownkey", false, slog.AnyValue(nil)),

		Entry("there is an attribute matching this group key",
			"request.method", true, slog.StringValue("POST")),

		Entry("there is no attribute matching this group key",
			"request.path", false, slog.AnyValue(nil)),
	)

	Describe("provides a getter for group names.", func() {

		BeforeEach(func() {
			logger = logger.WithGroup("g1").WithGroup("g2").WithGroup("g3")
			logger.Error("error message")
			var logFound bool
			log, logFound = observer.FindLog(1)
			Expect(logFound).To(BeTrue())
		})

		It("should return group names as a string, separated by dots", func() {
			Expect(log.GroupNames()).To(Equal("g1.g2.g3"))
		})
	})

})
