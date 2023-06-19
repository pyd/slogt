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

Note: no need to check that the slogt.ObserverHandler.Handle() method
is passing each slog.Record from the Logger to the observer as long as
the observer tests run as expected
*/

var _ = Describe("Handler", func() {

	// there are 2 constructors: NewObserverHandler(handler, observer) & NewDefaultObserverHandler(observer)
	Describe("Constructors", func() {

		var handlerArg slog.Handler
		var observerArg slogt.RecordCollector
		var constructorErr error

		BeforeEach(func() {
			// by default constructor arguments are not nil
			// nil argiumetns for testing will be set in sub BeforeEach's
			handlerArg = slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})
			observerArg = new(slogt.Observer)
		})

		Context("First constructor, with handler and observer arguments", func() {

			JustBeforeEach(func() {
				_, constructorErr = slogt.NewObserverHandler(handlerArg, observerArg)
			})

			When("handler argument is nil", func() {

				BeforeEach(func() {
					handlerArg = nil
				})

				It("should return an error", func() {
					Expect(constructorErr).To(HaveOccurred())
				})

			})

			When("observer argument is nil", func() {

				BeforeEach(func() {
					observerArg = nil
				})

				It("should return an error", func() {
					Expect(constructorErr).To(HaveOccurred())
				})

			})

			When("handler and observer argument are not nil", func() {

				It("should not return an error", func() {
					Expect(constructorErr).NotTo(HaveOccurred())
				})

			})
		})

		Context("Second constructor, with observer argument only (default handler)", func() {

			JustBeforeEach(func() {
				_, constructorErr = slogt.NewDefaultObserverHandler(observerArg)
			})

			When("observer argument is nil", func() {

				BeforeEach(func() {
					observerArg = nil
				})

				It("should return an error", func() {
					Expect(constructorErr).To(HaveOccurred())
				})

			})

			When("observer argument is not nil", func() {

				It("should not return an error", func() {
					Expect(constructorErr).NotTo(HaveOccurred())
				})

			})
		})

	})
})
