package slogt_test

import (
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/slog"

	"github.com/pyd/slogt"
)

var _ = Describe("Handler", func() {

	// var handler slog.Handler

	Describe("Constructors", func() {

		var handlerArg slog.Handler
		var observerArg slogt.RecordCollector
		var constructorErr error

		BeforeEach(func() {
			// set not nil args by default; overwrite with nil in sub BeforeEach
			handlerArg = slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})
			observerArg = new(slogt.Observer)
		})

		Context("With handler and observer arguments", func() {

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

		Context("With observer argument only (default handler)", func() {

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
