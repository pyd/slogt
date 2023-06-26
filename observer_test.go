package slogt_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/slog"

	"github.com/pyd/slogt"
)

/*
Testing slogt.Observer
*/

var _ = Describe("The log observer", func() {

	var observer *slogt.Observer
	var logger *slog.Logger

	BeforeEach(func() {
		observer = new(slogt.Observer)
		handler, _ := slogt.NewDefaultObserverHandler(observer)
		logger = slog.New(handler)
	})

	Describe("provides a log counter.", func() {

		When("no log was captured", func() {

			It("should return 0", func() {
				Expect(observer.CountLogs()).To(Equal(0))
			})

		})

		When("some logs where captured", func() {

			BeforeEach(func() {
				logger.Info("info message")
				logger.Warn("warn message")
				logger.Error("error message")
			})

			It("should return the exact number of logs", func() {
				Expect(observer.CountLogs()).To(Equal(3))
			})
		})
	})

	Describe("can find a log by it's chronological index", func() {

		var log slogt.Log
		var logFound bool

		BeforeEach(func() {
			logger.Warn("warn message 2")
			logger.Error("error message 2")
		})

		When("there is no log with this index", func() {

			BeforeEach(func() {
				log, logFound = observer.FindLog(99)
			})

			It("should return log not found", func() {
				Expect(logFound).To(BeFalse())
			})

			It("should return a zero-ed log", func() {
				Expect(log.Message()).To(BeEmpty())
			})

		})

		When("there is a log with this index", func() {

			BeforeEach(func() {
				log, logFound = observer.FindLog(2)
			})

			It("should return log found", func() {
				Expect(logFound).To(BeTrue())
			})

			It("should return the matching log", func() {
				Expect(log.Message()).To(Equal("error message 2"))
			})
		})
	})

	Describe("can clear all stored logs", func() {

		BeforeEach(func() {
			logger.Info("info message")
			logger.Warn("warn message")
			Expect(observer.CountLogs()).To(Equal(2))
		})

		It("should return a count of 0", func() {
			observer.ClearLogs()
			Expect(observer.CountLogs()).To(Equal(0))
		})
	})

	Describe("provides a getter for all logs", func() {

		var logs []slogt.Log

		JustBeforeEach(func() {
			logs = observer.Logs()
		})

		When("several logs were captured", func() {

			BeforeEach(func() {
				logger.Info("info message")
				logger.Warn("warn message")
			})

			It("should return the expected number of logs", func() {
				Expect(logs).To(HaveLen(2))
			})

			It("should return the expected logs", func() {
				log1, log1Found := observer.FindLog(1)
				Expect(log1Found).To(BeTrue())
				Expect(log1.Level()).To(Equal(slog.LevelInfo))

				log2, log2Found := observer.FindLog(2)
				Expect(log2Found).To(BeTrue())
				Expect(log2.Level()).To(Equal(slog.LevelWarn))
			})
		})

		When("no log was captured", func() {

			It("should return no logs", func() {
				Expect(logs).To(BeEmpty())
			})
		})
	})

})
