package client

// REQUEST

type YoutubeAccountRequest struct {
	Context            Context `json:"context"`
	DeviceTheme        string  `json:"deviceTheme"`
	UserInterfaceTheme string  `json:"userInterfaceTheme"`
}

type YoutubeRequest struct {
	Context  Context `json:"context"`
	BrowseID string  `json:"browseId"`
	Params   string  `json:"params"`
}

type ConfigInfo struct {
	AppInstallData string `json:"appInstallData"`
}

type MainAppWebInfo struct {
	GraftURL                  string `json:"graftUrl"`
	PwaInstallabilityStatus   string `json:"pwaInstallabilityStatus"`
	WebDisplayMode            string `json:"webDisplayMode"`
	IsWebNativeShareAvailable bool   `json:"isWebNativeShareAvailable"`
}

type client struct {
	Hl                 string         `json:"hl"`
	Gl                 string         `json:"gl"`
	RemoteHost         string         `json:"remoteHost"`
	DeviceMake         string         `json:"deviceMake"`
	DeviceModel        string         `json:"deviceModel"`
	VisitorData        string         `json:"visitorData"`
	UserAgent          string         `json:"userAgent"`
	ClientName         string         `json:"clientName"`
	ClientVersion      string         `json:"clientVersion"`
	OsName             string         `json:"osName"`
	OsVersion          string         `json:"osVersion"`
	OriginalURL        string         `json:"originalUrl"`
	Platform           string         `json:"platform"`
	ClientFormFactor   string         `json:"clientFormFactor"`
	ConfigInfo         ConfigInfo     `json:"configInfo"`
	UserInterfaceTheme string         `json:"userInterfaceTheme"`
	TimeZone           string         `json:"timeZone"`
	BrowserName        string         `json:"browserName"`
	BrowserVersion     string         `json:"browserVersion"`
	AcceptHeader       string         `json:"acceptHeader"`
	DeviceExperimentID string         `json:"deviceExperimentId"`
	ScreenWidthPoints  int            `json:"screenWidthPoints"`
	ScreenHeightPoints int            `json:"screenHeightPoints"`
	ScreenPixelDensity int            `json:"screenPixelDensity"`
	ScreenDensityFloat int            `json:"screenDensityFloat"`
	UtcOffsetMinutes   int            `json:"utcOffsetMinutes"`
	MemoryTotalKbytes  string         `json:"memoryTotalKbytes"`
	MainAppWebInfo     MainAppWebInfo `json:"mainAppWebInfo"`
}
type User struct {
	LockedSafetyMode bool `json:"lockedSafetyMode"`
}
type Request struct {
	UseSsl                  bool          `json:"useSsl"`
	InternalExperimentFlags []interface{} `json:"internalExperimentFlags"`
	ConsistencyTokenJars    []interface{} `json:"consistencyTokenJars"`
}
type ClickTracking struct {
	ClickTrackingParams string `json:"clickTrackingParams"`
}
type ParamsReq struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type AdSignalsInfo struct {
	Params []Params `json:"params"`
	Bid    string   `json:"bid"`
}

type Params struct {
}

type Context struct {
	Client        client        `json:"client"`
	User          User          `json:"user"`
	Request       Request       `json:"request"`
	ClickTracking ClickTracking `json:"clickTracking"`
	AdSignalsInfo AdSignalsInfo `json:"adSignalsInfo"`
}

var r = YoutubeRequest{
	Context: Context{
		AdSignalsInfo: AdSignalsInfo{
			Bid: "ANyPxKoknznTBcdHivD5_1ayWSNJ4NnHSSkoqIbSU5aGahJ8AbMTVbhu5WtSoCaKMbiueL70611eVWlXywKxMmiiLKoAxNVATw",
		},
		Client: client{
			AcceptHeader:     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			BrowserName:      "Chrome",
			BrowserVersion:   "107.0.0.0",
			ClientFormFactor: "UNKNOWN_FORM_FACTOR",
			ClientName:       "WEB",
			ClientVersion:    "2.20221107.06.00",
			DeviceMake:       "",
			DeviceModel:      "",
			Gl:               "US",
			Hl:               "en",
			MainAppWebInfo: MainAppWebInfo{
				GraftURL:                  "https://www.youtube.com/",
				PwaInstallabilityStatus:   "PWA_INSTALLABILITY_STATUS_CAN_BE_INSTALLED",
				WebDisplayMode:            "WEB_DISPLAY_MODE_BROWSER",
				IsWebNativeShareAvailable: true,
			},
			MemoryTotalKbytes:  "8000000",
			OriginalURL:        "https://www.youtube.com/",
			OsName:             "Windows",
			OsVersion:          "10.0",
			Platform:           "DESKTOP",
			ScreenDensityFloat: 1,
			ScreenHeightPoints: 969,
			ScreenPixelDensity: 1,
			ScreenWidthPoints:  659,
			UserAgent:          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36,gzip(gfe)",
			UserInterfaceTheme: "USER_INTERFACE_THEME_DARK",
		},
		Request: Request{
			ConsistencyTokenJars:    []interface{}{},
			InternalExperimentFlags: []interface{}{},
			UseSsl:                  true,
		},
		ClickTracking: ClickTracking{
			ClickTrackingParams: "COIBEMCeCRgAIhMI7JOCuYTq-wIV8sxPCB3rjQKr",
		},
		User: User{
			LockedSafetyMode: false,
		},
	},
	Params: "EgVhYm91dPIGBAoCEgA%3D",
}

// RESPONSE

type YoutubeAccountData struct {
	Code string `json:"code"`
	Data struct {
		ResponseContext struct {
			ServiceTrackingParams []struct {
				Service string `json:"service"`
				Params  []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"params"`
			} `json:"serviceTrackingParams"`
			MainAppWebResponseContext struct {
				DatasyncID string `json:"datasyncId"`
				LoggedOut  bool   `json:"loggedOut"`
			} `json:"mainAppWebResponseContext"`
			WebResponseContextExtensionData struct {
				HasDecorated bool `json:"hasDecorated"`
			} `json:"webResponseContextExtensionData"`
		} `json:"responseContext"`
		SelectText struct {
			SimpleText string `json:"simpleText"`
		} `json:"selectText"`
		Actions []struct {
			GetMultiPageMenuAction struct {
				Menu struct {
					MultiPageMenuRenderer struct {
						Header struct {
							SimpleMenuHeaderRenderer struct {
								BackButton struct {
									ButtonRenderer struct {
										Style      string `json:"style"`
										Size       string `json:"size"`
										IsDisabled bool   `json:"isDisabled"`
										Icon       struct {
											IconType string `json:"iconType"`
										} `json:"icon"`
										Accessibility struct {
											Label string `json:"label"`
										} `json:"accessibility"`
										AccessibilityData struct {
											AccessibilityData struct {
												Label string `json:"label"`
											} `json:"accessibilityData"`
										} `json:"accessibilityData"`
									} `json:"buttonRenderer"`
								} `json:"backButton"`
								Title struct {
									SimpleText string `json:"simpleText"`
								} `json:"title"`
							} `json:"simpleMenuHeaderRenderer"`
						} `json:"header"`
						Sections []struct {
							AccountSectionListRenderer struct {
								Contents []struct {
									AccountItemSectionRenderer struct {
										Contents []struct {
											AccountItem struct {
												AccountName struct {
													SimpleText string `json:"simpleText"`
												} `json:"accountName"`
												AccountPhoto struct {
													Thumbnails []struct {
														URL    string `json:"url"`
														Width  int    `json:"width"`
														Height int    `json:"height"`
													} `json:"thumbnails"`
												} `json:"accountPhoto"`
												IsSelected   bool `json:"isSelected"`
												IsDisabled   bool `json:"isDisabled"`
												MobileBanner struct {
													Thumbnails []struct {
														URL    string `json:"url"`
														Width  int    `json:"width"`
														Height int    `json:"height"`
													} `json:"thumbnails"`
												} `json:"mobileBanner"`
												HasChannel      bool `json:"hasChannel"`
												ServiceEndpoint struct {
													SelectActiveIdentityEndpoint struct {
														SupportedTokens []struct {
															AccountStateToken struct {
																HasChannel       bool   `json:"hasChannel"`
																IsMerged         bool   `json:"isMerged"`
																ObfuscatedGaiaID string `json:"obfuscatedGaiaId"`
															} `json:"accountStateToken,omitempty"`
															OfflineCacheKeyToken struct {
																ClientCacheKey string `json:"clientCacheKey"`
															} `json:"offlineCacheKeyToken,omitempty"`
															AccountSigninToken struct {
																SigninURL string `json:"signinUrl"`
															} `json:"accountSigninToken,omitempty"`
															DatasyncIDToken struct {
																DatasyncIDToken string `json:"datasyncIdToken"`
															} `json:"datasyncIdToken,omitempty"`
														} `json:"supportedTokens"`
														NextNavigationEndpoint struct {
															CommandMetadata struct {
																WebCommandMetadata struct {
																	URL         string `json:"url"`
																	WebPageType string `json:"webPageType"`
																	RootVe      int    `json:"rootVe"`
																} `json:"webCommandMetadata"`
															} `json:"commandMetadata"`
															URLEndpoint struct {
																URL string `json:"url"`
															} `json:"urlEndpoint"`
														} `json:"nextNavigationEndpoint"`
													} `json:"selectActiveIdentityEndpoint"`
												} `json:"serviceEndpoint"`
												AccountByline struct {
													SimpleText string `json:"simpleText"`
												} `json:"accountByline"`
											} `json:"accountItem,omitempty"`
											CompactLinkRenderer struct {
												Title struct {
													SimpleText string `json:"simpleText"`
												} `json:"title"`
												NavigationEndpoint struct {
													CommandMetadata struct {
														WebCommandMetadata struct {
															URL         string `json:"url"`
															WebPageType string `json:"webPageType"`
															RootVe      int    `json:"rootVe"`
														} `json:"webCommandMetadata"`
													} `json:"commandMetadata"`
													SignalNavigationEndpoint struct {
														Signal string `json:"signal"`
													} `json:"signalNavigationEndpoint"`
												} `json:"navigationEndpoint"`
											} `json:"compactLinkRenderer,omitempty"`
										} `json:"contents"`
									} `json:"accountItemSectionRenderer"`
								} `json:"contents"`
								Header struct {
									GoogleAccountHeaderRenderer struct {
										Name struct {
											SimpleText string `json:"simpleText"`
										} `json:"name"`
										Email struct {
											SimpleText string `json:"simpleText"`
										} `json:"email"`
									} `json:"googleAccountHeaderRenderer"`
								} `json:"header"`
							} `json:"accountSectionListRenderer"`
						} `json:"sections"`
						Footer struct {
							MultiPageMenuSectionRenderer struct {
								Items []struct {
									CompactLinkRenderer struct {
										Icon struct {
											IconType string `json:"iconType"`
										} `json:"icon"`
										Title struct {
											SimpleText string `json:"simpleText"`
										} `json:"title"`
										NavigationEndpoint struct {
											CommandMetadata struct {
												WebCommandMetadata struct {
													URL         string `json:"url"`
													WebPageType string `json:"webPageType"`
													RootVe      int    `json:"rootVe"`
												} `json:"webCommandMetadata"`
											} `json:"commandMetadata"`
											URLEndpoint struct {
												URL string `json:"url"`
											} `json:"urlEndpoint"`
										} `json:"navigationEndpoint"`
										Style string `json:"style"`
									} `json:"compactLinkRenderer"`
								} `json:"items"`
							} `json:"multiPageMenuSectionRenderer"`
						} `json:"footer"`
						Style string `json:"style"`
					} `json:"multiPageMenuRenderer"`
				} `json:"menu"`
			} `json:"getMultiPageMenuAction"`
		} `json:"actions"`
	} `json:"data"`
}

// ChannelInfo Response

type YoutubeChannelData struct {
	Contents        Contents        `json:"contents"`
	Header          Header          `json:"header"`
	ResponseContext ResponseContext `json:"responseContext"`
}

type Contents struct {
	TwoColumnBrowseResultsRenderer TwoColumnBrowseResultsRenderer `json:"twoColumnBrowseResultsRenderer"`
}

type TwoColumnBrowseResultsRenderer struct {
	Tabs []Tab `json:"tabs"`
}

type Tab struct {
	TabRenderer TabRenderer `json:"tabRenderer"`
}

type TabRenderer struct {
	RendererContent RendererContent `json:"content"`
	Title           string          `json:"title"`
}

type RendererContent struct {
	SectionListRenderer SectionListRenderer `json:"sectionListRenderer"`
}

type SectionListRenderer struct {
	ListRendererContents []ListRendererContent `json:"contents"`
	DisablePullToRefresh bool                  `json:"disablePullToRefresh"`
}

type ListRendererContent struct {
	ItemSectionRenderer ItemSectionRenderer `json:"itemSectionRenderer"`
}

type ItemSectionRenderer struct {
	SectionContents []SectionContent `json:"contents"`
	Tracking        string           `json:"trackingParams"`
}

type SectionContent struct {
	ChannelAboutFullMetadataRenderer ChannelAboutFullMetadataRenderer `json:"channelAboutFullMetadataRenderer"`
}

type ChannelAboutFullMetadataRenderer struct {
	ChannelId      string         `json:"channelId"`
	ViewCountText  ViewCountText  `json:"viewCountText"`
	JoinedDateText JoinedDateText `json:"joinedDateText"`
	Title          Title          `json:"title"`
}

type ViewCountText struct {
	SimpleText string `json:"simpleText"`
}

type Title struct {
	SimpleText string `json:"simpleText"`
}

type Header struct {
	C4TabbedHeaderRenderer C4TabbedHeaderRenderer `json:"c4TabbedHeaderRenderer"`
}

type C4TabbedHeaderRenderer struct {
	Title               string              `json:"title"`
	Badges              []Badges            `json:"badges"`
	SubscriberCountText SubscriberCountText `json:"subscriberCountText"`
	VideosCountText     VideosCountText     `json:"videosCountText"`
}

type Badges struct {
	MetaDataBadgeRenderer MetaDataBadgeRenderer `json:"metadataBadgeRenderer"`
}

type MetaDataBadgeRenderer struct {
	Style string `json:"style"`
}

type SubscriberCountText struct {
	SimpleText string `json:"simpleText"`
}

type VideosCountText struct {
	Runs []VideoRun `json:"runs"`
}

type JoinedDateText struct {
	Runs []DateRun `json:"runs"`
}

type VideoRun struct {
	Text string `json:"text"`
}

type DateRun struct {
	Text string `json:"text"`
}

// ///////////////////// API ////////////////////////
type YoutubeAPIChannelData struct {
	Items []ApiItem `json:"items"`
}

type ApiItem struct {
	Snippet    Snippet    `json:"snippet"`
	Statistics Statistics `json:"statistics"`
}

type Snippet struct {
	RegDate string `json:"publishedAt"`
	Country string `json:"country"`
}

type Statistics struct {
	ViewCount   string `json:"viewCount"`
	Subscribers string `json:"subscriberCount"`
	VideoCount  string `json:"videoCount"`
}

type ResponseContext struct {
	ServiceTrackingParams []ServiceTrackingParams `json:"serviceTrackingParams"`
}

type ServiceTrackingParams struct {
	Service        string           `json:"service"`
	TrackingParams []TrackingParams `json:"trackingParams"`
}

type TrackingParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
