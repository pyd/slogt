package slogt_test

import (
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/slog"

	"github.com/pyd/slogt"
)

/*
Testing slogt.ObserverHandler

Check that the Handle() method is passing each slog.Record, groups and shared attributes to the observer is implicit: see Observer tests
Check that the WithAttrs() method captures shared attributes is implicit: see Log tests
Check that the WithGroup() method captures groups is implicit: see Log tests
*/

var _ = Describe("Observer Handler", func() {

	// 2 constructors: NewObserverHandler(observer, handler) & NewDefaultObserverHandler(observer)
	Describe("has 2 constructors.", func() {

		var handlerArg slog.Handler
		var observerArg slogt.HandlerObserver
		var constructorErr error

		BeforeEach(func() {
			// by default constructor arguments are not nil
			// nil arguments for testing will be set in sub BeforeEach's
			handlerArg = slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})
			observerArg = new(slogt.Observer)
		})

		Context("One requires a handler and an observer arguments", func() {

			JustBeforeEach(func() {
				_, constructorErr = slogt.NewObserverHandler(observerArg, handlerArg)
			})

			When("the handler argument is nil", func() {

				BeforeEach(func() {
					handlerArg = nil
				})

				It("should return an error", func() {
					Expect(constructorErr).To(HaveOccurred())
				})

			})

			When("the observer argument is nil", func() {

				BeforeEach(func() {
					observerArg = nil
				})

				It("should return an error", func() {
					Expect(constructorErr).To(HaveOccurred())
				})

			})

			When("the handler and observer arguments are not nil", func() {

				It("should not return an error", func() {
					Expect(constructorErr).NotTo(HaveOccurred())
				})

			})
		})

		Context("One requires only an observer argument", func() {

			JustBeforeEach(func() {
				_, constructorErr = slogt.NewDefaultObserverHandler(observerArg)
			})

			When("this argument is nil", func() {

				BeforeEach(func() {
					observerArg = nil
				})

				It("should return an error", func() {
					Expect(constructorErr).To(HaveOccurred())
				})

			})

			When("this argument is not nil", func() {

				It("should not return an error", func() {
					Expect(constructorErr).NotTo(HaveOccurred())
				})

			})
		})

	})

	Describe("provides a getter for groups", func() {

		var handler slogt.ObserverHandler

		BeforeEach(func() {
			// create a handler with group "app1"
			observer := new(slogt.Observer)
			handler, _ = slogt.NewDefaultObserverHandler(observer)
			handler = handler.WithGroup("app1").(slogt.ObserverHandler)
			// create a logger with this handler and add group "admin"
			logger := slog.New(handler)
			logger = logger.WithGroup("admin") // add group to the logger handler
			// we need to get the handler from the logger to see the "admin" group
			handler = logger.Handler().(slogt.ObserverHandler)
		})

		It("should return groups from handler and logger", func() {
			Expect(handler.Groups()).To(HaveExactElements("app1", "admin"))
		})
	})

	Describe("provides a getter and a finder for attributes.", func() {

		var handler slogt.ObserverHandler
		var attributes []slog.Attr
		// result of attribute finder
		var attribute slog.Attr
		var attributeFound bool

		BeforeEach(func() {
			// create a handler with group "app1" and an attribute
			observer := new(slogt.Observer)
			handler, _ = slogt.NewDefaultObserverHandler(observer)
			handler = handler.WithGroup("app1").(slogt.ObserverHandler)
			handler = handler.WithAttrs([]slog.Attr{slog.String("handler-attr", "handler attr value")}).(slogt.ObserverHandler)
			// create a logger with this handler and add group "admin"
			logger := slog.New(handler)
			logger = logger.WithGroup("admin")
			logger = logger.With(slog.String("logger-attr", "logger attr value"))
			// get back handler from the logger
			handler = logger.Handler().(slogt.ObserverHandler)
			// check attributes getter
			attributes = handler.Atttributes()
			// check attributes finder
			attribute, attributeFound = handler.FindAttribute("logger-attr")
		})

		Describe("The getter", func() {

			It("should return 2 attributes", func() {
				Expect(attributes).To(HaveLen(2))
			})

			It("should return the handler attribute", func() {
				Expect(attributes[0].Key).To(Equal("handler-attr"))
				Expect(attributes[0].Value.String()).To(Equal("handler attr value"))
			})

			It("should return the logger attribute", func() {
				Expect(attributes[1].Key).To(Equal("logger-attr"))
				Expect(attributes[1].Value.String()).To(Equal("logger attr value"))
			})
		})

		Describe("The finder by Key", func() {

			When("the key exists", func() {

				It("should return attribute found", func() {
					Expect(attributeFound).To(BeTrue())
				})

				It("should return the matching attribute", func() {
					Expect(attribute.Value.String()).To(Equal("logger attr value"))
				})
			})

			When("the key does not exist", func() {

				BeforeEach(func() {
					attribute, attributeFound = handler.FindAttribute("unknownkey")
				})

				It("should return attribute not found", func() {
					Expect(attributeFound).To(BeFalse())
				})

				It("should return an attribute with a nil Value", func() {
					Expect(attribute.Value.Any()).To(BeNil())
				})
			})

		})
	})

})
