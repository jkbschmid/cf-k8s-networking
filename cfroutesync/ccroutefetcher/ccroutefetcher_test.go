package ccroutefetcher_test

import (
	"errors"

	"code.cloudfoundry.org/cf-k8s-networking/cfroutesync/ccclient"
	"code.cloudfoundry.org/cf-k8s-networking/cfroutesync/ccroutefetcher"
	"code.cloudfoundry.org/cf-k8s-networking/cfroutesync/ccroutefetcher/fakes"
	"code.cloudfoundry.org/cf-k8s-networking/cfroutesync/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetching once", func() {
	var (
		fakeCCClient     *fakes.CCClient
		fakeUAAClient    *fakes.UAAClient
		fakeSnapshotRepo *fakes.SnapshotRepo
		expectedSnapshot *models.RouteSnapshot
		fetcher          *ccroutefetcher.Fetcher
		routesList       []ccclient.Route
	)

	BeforeEach(func() {
		fakeRoute0Destination0 := ccclient.Destination{
			Guid:   "route-0-dest-0-guid",
			Weight: models.IntPtr(10),
			Port:   8000,
		}
		fakeRoute0Destination0.App.Guid = "route-0-dest-0-app-0-guid"
		fakeRoute0Destination0.App.Process.Type = "route-0-dest-0-app-0-process-type"

		fakeRoute0Destination1 := ccclient.Destination{
			Guid:   "route-0-dest-1-guid",
			Weight: models.IntPtr(11),
			Port:   8001,
		}
		fakeRoute0Destination1.App.Guid = "route-0-dest-1-app-1-guid"
		fakeRoute0Destination1.App.Process.Type = "route-0-dest-1-app-1-process-type"

		fakeRoute1Destination0 := ccclient.Destination{
			Guid:   "route-1-dest-0-guid",
			Weight: models.IntPtr(12),
			Port:   9000,
		}
		fakeRoute1Destination0.App.Guid = "route-1-dest-0-app-0-guid"
		fakeRoute1Destination0.App.Process.Type = "route-1-dest-0-app-0-process-type"

		fakeCCClient = &fakes.CCClient{}

		routesList = []ccclient.Route{
			ccclient.Route{
				Guid: "route-0-guid",
				Host: "route-0-host",
				Path: "/route-0-path",
				Url:  "route-0-host.domain0.example.com/route-0-path",
				Destinations: []ccclient.Destination{
					fakeRoute0Destination0,
					fakeRoute0Destination1,
				},
			},
			ccclient.Route{
				Guid: "route-1-guid",
				Host: "route-1-host",
				Path: "/route-1-path",
				Url:  "route-1-host.domain1.apps.internal/route-1-path",
				Destinations: []ccclient.Destination{
					fakeRoute1Destination0,
				},
			},
			ccclient.Route{
				Guid:         "route-2-guid",
				Host:         "route-2-host",
				Path:         "/route-2-path",
				Url:          "route-2-host.domain1.apps.internal/route-2-path",
				Destinations: []ccclient.Destination{},
			},
		}
		routesList[0].Relationships.Domain.Data.Guid = "domain-0-guid"
		routesList[1].Relationships.Domain.Data.Guid = "domain-1-guid"
		routesList[2].Relationships.Domain.Data.Guid = "domain-1-guid"
		routesList[0].Relationships.Space.Data.Guid = "space-0-guid"
		routesList[1].Relationships.Space.Data.Guid = "space-1-guid"
		routesList[2].Relationships.Space.Data.Guid = "space-1-guid"
		fakeCCClient.ListRoutesReturns(routesList, nil)

		fakeCCClient.ListDomainsReturns([]ccclient.Domain{
			{
				Guid:     "domain-0-guid",
				Name:     "domain0.example.com",
				Internal: false,
			},
			{
				Guid:     "domain-1-guid",
				Name:     "domain1.apps.internal",
				Internal: true,
			},
		}, nil)

		spacesList := []ccclient.Space{
			{
				Guid: "space-0-guid",
			},
			{
				Guid: "space-1-guid",
			},
		}
		spacesList[0].Relationships.Organization.Data.Guid = "org-0-guid"
		spacesList[1].Relationships.Organization.Data.Guid = "org-1-guid"
		fakeCCClient.ListSpacesReturns(spacesList, nil)

		fakeUAAClient = &fakes.UAAClient{}
		fakeUAAClient.GetTokenReturns("fake-uaa-token", nil)

		fakeSnapshotRepo = &fakes.SnapshotRepo{}

		expectedSnapshot = &models.RouteSnapshot{
			Routes: []models.Route{
				models.Route{
					Guid: "route-0-guid",
					Host: "route-0-host",
					Path: "/route-0-path",
					Url:  "route-0-host.domain0.example.com/route-0-path",
					Domain: models.Domain{
						Guid:     "domain-0-guid",
						Name:     "domain0.example.com",
						Internal: false,
					},
					Space: models.Space{
						Guid: "space-0-guid",
						Organization: models.Organization{
							Guid: "org-0-guid",
						},
					},
					Destinations: []models.Destination{
						models.Destination{
							Guid: "route-0-dest-0-guid",
							App: models.App{
								Guid:    "route-0-dest-0-app-0-guid",
								Process: models.Process{Type: "route-0-dest-0-app-0-process-type"},
							},
							Port:   8000,
							Weight: models.IntPtr(10),
						},
						models.Destination{
							Guid: "route-0-dest-1-guid",
							App: models.App{
								Guid:    "route-0-dest-1-app-1-guid",
								Process: models.Process{Type: "route-0-dest-1-app-1-process-type"},
							},
							Port:   8001,
							Weight: models.IntPtr(11),
						},
					},
				},
				models.Route{
					Guid: "route-1-guid",
					Host: "route-1-host",
					Path: "/route-1-path",
					Url:  "route-1-host.domain1.apps.internal/route-1-path",
					Domain: models.Domain{
						Guid:     "domain-1-guid",
						Name:     "domain1.apps.internal",
						Internal: true,
					},
					Space: models.Space{
						Guid: "space-1-guid",
						Organization: models.Organization{
							Guid: "org-1-guid",
						},
					},
					Destinations: []models.Destination{
						models.Destination{
							Guid: "route-1-dest-0-guid",
							App: models.App{
								Guid:    "route-1-dest-0-app-0-guid",
								Process: models.Process{Type: "route-1-dest-0-app-0-process-type"},
							},
							Port:   9000,
							Weight: models.IntPtr(12),
						},
					},
				},
				models.Route{
					Guid: "route-2-guid",
					Host: "route-2-host",
					Path: "/route-2-path",
					Url:  "route-2-host.domain1.apps.internal/route-2-path",
					Domain: models.Domain{
						Guid:     "domain-1-guid",
						Name:     "domain1.apps.internal",
						Internal: true,
					},
					Space: models.Space{
						Guid: "space-1-guid",
						Organization: models.Organization{
							Guid: "org-1-guid",
						},
					},
					Destinations: nil,
				},
			},
		}

		fetcher = &ccroutefetcher.Fetcher{
			CCClient:     fakeCCClient,
			UAAClient:    fakeUAAClient,
			SnapshotRepo: fakeSnapshotRepo,
		}
	})

	It("calls cc client to get routes and destinations", func() {
		err := fetcher.FetchOnce()
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeCCClient.ListRoutesCallCount()).To(Equal(1))
		token := fakeCCClient.ListRoutesArgsForCall(0)
		Expect(token).To(Equal("fake-uaa-token"))
	})

	Context("when there are routes to save", func() {
		It("converts cc types to a route snapshot and puts that into the repo", func() {

			err := fetcher.FetchOnce()
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeSnapshotRepo.PutCallCount()).To(Equal(1))
			Expect(fakeSnapshotRepo.PutArgsForCall(0)).To(Equal(expectedSnapshot))
		})

		It("transforms host and domain to lower case to have DNS-1123 compliant fqdns", func() {
			routesList = []ccclient.Route{
				{
					Guid: "uppercase-route-guid",
					Host: "UPPERCASE-HOST",
					Path: "/PATH",
					Url:  "UPPERCASE-HOST.EXAMPLE.COM/PATH",
				},
			}
			routesList[0].Relationships.Domain.Data.Guid = "domain-0-guid"
			routesList[0].Relationships.Space.Data.Guid = "space-0-guid"
			fakeCCClient.ListRoutesReturns(routesList, nil)

			fakeCCClient.ListDomainsReturns([]ccclient.Domain{
				{
					Guid:     "domain-0-guid",
					Name:     "EXAMPLE.COM",
					Internal: false,
				},
			}, nil)

			expectedSnapshot = &models.RouteSnapshot{
				Routes: []models.Route{
					models.Route{
						Guid: "uppercase-route-guid",
						Host: "uppercase-host",
						Path: "/PATH",
						Url:  "uppercase-host.example.com/PATH",
						Domain: models.Domain{
							Guid:     "domain-0-guid",
							Name:     "example.com",
							Internal: false,
						},
						Space: models.Space{
							Guid:         "space-0-guid",
							Organization: models.Organization{Guid: "org-0-guid"},
						},
						Destinations: nil,
					},
				},
			}

			err := fetcher.FetchOnce()
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeSnapshotRepo.PutCallCount()).To(Equal(1))
			Expect(fakeSnapshotRepo.PutArgsForCall(0)).To(Equal(expectedSnapshot))
		})
	})

	Context("when there is an error getting the token from UAA", func() {
		It("returns the error", func() {
			fakeUAAClient.GetTokenReturns("", errors.New("banana"))
			err := fetcher.FetchOnce()
			Expect(err).To(MatchError("uaa get token: banana"))
		})
	})

	Context("when there is an error getting Routes from Cloud Controller", func() {
		It("returns the error", func() {
			fakeCCClient.ListRoutesReturns(nil, errors.New("potato!"))
			err := fetcher.FetchOnce()
			Expect(err).To(MatchError("cc list routes: potato!"))
		})
	})

	Context("when there is an error getting Domains from Cloud Controller", func() {
		It("returns the error", func() {
			fakeCCClient.ListDomainsReturns(nil, errors.New("ohno!"))
			err := fetcher.FetchOnce()
			Expect(err).To(MatchError("cc list domains: ohno!"))
		})
	})

	Context("when a route refers to a domain that was not found", func() {
		It("returns an error", func() {
			fakeCCClient.ListDomainsReturns([]ccclient.Domain{
				{
					Guid:     "domain-1-guid",
					Name:     "domain1.apps.internal",
					Internal: true,
				},
			}, nil)
			err := fetcher.FetchOnce()
			Expect(err).To(MatchError("route route-0-guid refers to missing domain domain-0-guid"))
		})
	})

	Context("when there is an error getting Spaces from Cloud Controller", func() {
		It("returns the error", func() {
			fakeCCClient.ListSpacesReturns(nil, errors.New("ohno!"))
			err := fetcher.FetchOnce()
			Expect(err).To(MatchError("cc list spaces: ohno!"))
		})
	})

	Context("when a route refers to a Space that was not found", func() {
		It("returns an error", func() {
			fakeCCClient.ListSpacesReturns([]ccclient.Space{
				{
					Guid: "not-space-0",
				},
			}, nil)
			err := fetcher.FetchOnce()
			Expect(err).To(MatchError("route route-0-guid refers to missing space space-0-guid"))
		})
	})
})
