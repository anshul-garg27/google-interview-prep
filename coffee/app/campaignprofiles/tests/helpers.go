package tests

import (
	discoverydao "coffee/app/discovery/dao"
	beatservice "coffee/client/beat"
	coredomain "coffee/core/domain"
)

func prepareIaProfileForTest() *coredomain.Profile {
	code := "IA_2624407"
	platformCode := "2624407"
	name := "Virat Kohli"
	handle := "virat.kohli"
	username := "virat.kohli"
	igId := "2094200507"
	description := "Carpediem!"
	profileType := "Mega"
	thumbnail := "https://d24w28i6lzk071.cloudfront.net/assets/profiles/instagram/2094200507.jpg"
	followers := int64(259316402)
	storyReach := int64(92493)
	imageReach := int64(12207087)
	reelsReach := int64(17650056)
	avgReach := float64(13915888.561930226)
	avgReelsPlayCount := float64(19352073.916666668)
	avgLikes := float64(2540514)
	avgComments := float64(21929.083333333332)
	avgViews := float64(1874125.25)
	following := int64(286)
	return &coredomain.Profile{
		Code:         &code,
		Locked:       false,
		GccCode:      &code,
		PlatformCode: platformCode,
		Name:         &name,
		Handle:       &handle,
		Username:     &username,
		IgId:         &igId,
		Description:  &description,
		Platform:     "INSTAGRAM",
		ProfileType:  &profileType,
		Thumbnail:    &thumbnail,
		Metrics: coredomain.ProfileMetrics{
			Followers:                &followers,
			EngagementRatePercenatge: 0.9881531070037496,
			StoryReach:               &storyReach,
			ImageReach:               &imageReach,
			ReelsReach:               &reelsReach,
			AvgReach:                 &avgReach,
			AvgReelsPlayCount:        &avgReelsPlayCount,
			AvgLikes:                 &avgLikes,
			AvgComments:              &avgComments,
			AvgViews:                 &avgViews,
			Following:                &following,
		},
	}
}

func mockIaBeatResponse() *beatservice.CompleteProfileResponse {
	staticAvgReach := float64(12671464.125)
	reelsAvgPlays := float64(15933665.125)
	reelsAvgReach := float64(14532293.875)
	storyAvgReach := float64(5225562)
	allAvgER := float64(0.5445210507630777)
	allAvgLikes := float64(1400624.625)
	allAvgComments := float64(11572.25)
	return &beatservice.CompleteProfileResponse{
		CompleteProfile: beatservice.CompleteProfile{
			Platform:      "instagram",
			Handle:        "virat.kohli",
			ProfilePicUrl: "https://scontent-bom1-2.xx.fbcdn.net/v/t51.2885-15/331017149_3373809859551320_1963035851400324431_n.jpg?_nc_cat=1\\u0026ccb=1-7\\u0026_nc_sid=86c713\\u0026_nc_ohc=bQe6YXCRBGIAX__s4Nq\\u0026_nc_ht=scontent-bom1-2.xx\\u0026edm=AL-3X8kEAAAA\\u0026oh=00_AfAaeZvYvBlFctQqfcs5KwS5U4o2cCF2zRvBDHc2RrM93g\\u0026oe=65105D3E",
			ProfileId:     "2094200507",
			Biography:     "Carpediem!",
			Following:     286,
			Followers:     259346608,
			FullName:      "Virat Kohli",
			PostMetrics: beatservice.CompleteProfilePostMetrics{
				Static: beatservice.PostMetrics{
					AverageReach: &staticAvgReach,
				},
				Reels: beatservice.PostMetrics{
					AveragePlays: &reelsAvgPlays,
					AverageReach: &reelsAvgReach,
				},
				Story: beatservice.PostMetrics{
					AverageReach: &storyAvgReach,
				},
				All: beatservice.PostMetrics{
					AverageEngagementRate: &allAvgER,
					AverageLikes:          &allAvgLikes,
					AverageComments:       &allAvgComments,
				},
			},
		},
	}
}

func mockIaEntityToCheckBeatResponse() discoverydao.InstagramAccountEntity {
	expectedFullName := "Virat Kohli"
	expectedHandle := "virat.kohli"
	expectedIgId := "2094200507"
	expectedThumbnail := "https://scontent-bom1-2.xx.fbcdn.net/v/t51.2885-15/331017149_3373809859551320_1963035851400324431_n.jpg?_nc_cat=1\\u0026ccb=1-7\\u0026_nc_sid=86c713\\u0026_nc_ohc=bQe6YXCRBGIAX__s4Nq\\u0026_nc_ht=scontent-bom1-2.xx\\u0026edm=AL-3X8kEAAAA\\u0026oh=00_AfAaeZvYvBlFctQqfcs5KwS5U4o2cCF2zRvBDHc2RrM93g\\u0026oe=65105D3E"
	expectedBio := "Carpediem!"
	expectedFollowers := int64(259346608)
	expectedFollowing := int64(286)
	expectedAvgLikes := float64(1400624.625)
	expectedAvgReach := float64(13915888.561930226)
	expectedReelsAvgReach := int64(14532293)
	expectedStoryAvgReach := int64(5225562)
	expectedImageAvgReach := int64(12671464)
	expectedAvgComments := float64(11572.25)
	expectedAvgViews := float64(1874125.25)
	expectedReelsAvgPlayCount := float64(15933665.125)
	label := "Mega"
	profileId := "IA_2624407"
	return discoverydao.InstagramAccountEntity{
		ID:                2624407,
		ProfileId:         &profileId,
		Name:              &expectedFullName,
		Handle:            &expectedHandle,
		IgID:              &expectedIgId,
		Thumbnail:         &expectedThumbnail,
		Bio:               &expectedBio,
		Followers:         &expectedFollowers,
		Following:         &expectedFollowing,
		EngagementRate:    0.005445210507630777,
		AvgLikes:          &expectedAvgLikes,
		AvgReach:          &expectedAvgReach,
		ReelsReach:        &expectedReelsAvgReach,
		StoryReach:        &expectedStoryAvgReach,
		ImageReach:        &expectedImageAvgReach,
		AvgComments:       &expectedAvgComments,
		AvgViews:          &expectedAvgViews,
		AvgReelsPlayCount: &expectedReelsAvgPlayCount,
		Label:             &label,
	}
}

func prepareIaCampaignProfileForTest() *coredomain.CampaignProfileEntry {
	profileCode := "IA_1351727"
	name := "Virat Kohli"
	handle := "virat.kohli"
	username := "virat.kohli"
	thumbnail := "https://scontent-sea1-1.cdninstagram.com/v/t51.2885-19/331017149_3373809859551320_1963035851400324431_n.jpg?_nc_ht=scontent-sea1-1.cdninstagram.com\\u0026_nc_cat=1\\u0026_nc_ohc=c6MJT5zJGkwAX_ICnxD\\u0026edm=AEF8tYYBAAAA\\u0026ccb=7-5\\u0026oh=00_AfArTglQQ0hLiUn9T7YKdQjAq2xbEJJimEMmvKxlGmJj_g\\u0026oe=65113356\\u0026_nc_sid=1e20d2"
	followers := int64(259354757)
	following := int64(286)
	avgLikes := float64(1407920.875)
	avgComments := float64(11721.75)
	uploads := int64(1651)
	avgReelsPlay30d := int64(60946411)
	avgReach := float64(15687749.09291251)
	storyReach := int64(5184299)
	imageReach := int64(12673487)
	reelsReach := int64(14657218)
	isVerified := false
	accountType := "CREATOR"
	audienceCityAvailable := true
	audienceGenderAvailable := false
	igId := "2094200507"
	return &coredomain.CampaignProfileEntry{
		SocialDetails: &coredomain.SocialAccount{
			IgId:                     &igId,
			Code:                     int64(1351727),
			ProfileCode:              &profileCode,
			GCCProfileCode:           &profileCode,
			Name:                     &name,
			Platform:                 "INSTAGRAM",
			Handle:                   &handle,
			Username:                 &username,
			PlatformThumbnail:        &thumbnail,
			EngagementRatePercenatge: 0.5473748164179614,
			Followers:                &followers,
			Following:                &following,
			AvgLikes:                 &avgLikes,
			AvgComments:              &avgComments,
			Uploads:                  &uploads,
			AvgReelsPlay30d:          &avgReelsPlay30d,
			AvgReach:                 &avgReach,
			StoryReach:               &storyReach,
			ImageReach:               &imageReach,
			ReelsReach:               &reelsReach,
			IsVerified:               &isVerified,
			Source:                   "GCC",
			AccountType:              &accountType,
			AudienceCityAvailable:    &audienceCityAvailable,
			AudienceGenderAvailable:  &audienceGenderAvailable,
		},
	}
}

func mockIaCpEntryToCheckBeatResponse() *coredomain.CampaignProfileEntry {
	profileCode := "IA_1351727"
	name := "Virat Kohli"
	handle := "virat.kohli"
	username := "virat.kohli"
	thumbnail := "https://scontent-bom1-2.xx.fbcdn.net/v/t51.2885-15/331017149_3373809859551320_1963035851400324431_n.jpg?_nc_cat=1\\u0026ccb=1-7\\u0026_nc_sid=86c713\\u0026_nc_ohc=bQe6YXCRBGIAX__s4Nq\\u0026_nc_ht=scontent-bom1-2.xx\\u0026edm=AL-3X8kEAAAA\\u0026oh=00_AfAaeZvYvBlFctQqfcs5KwS5U4o2cCF2zRvBDHc2RrM93g\\u0026oe=65105D3E"
	followers := int64(259346608)
	following := int64(286)
	avgLikes := float64(1400624.625)
	avgComments := float64(11572.25)
	uploads := int64(1651)
	avgReelsPlay30d := int64(15933665)
	avgReach := float64(15687749.09291251)
	storyReach := int64(5225562)
	imageReach := int64(12671464)
	reelsReach := int64(14532293)
	isVerified := false
	accountType := "CREATOR"
	audienceCityAvailable := true
	audienceGenderAvailable := false
	igId := "2094200507"
	return &coredomain.CampaignProfileEntry{
		SocialDetails: &coredomain.SocialAccount{
			IgId:                     &igId,
			Code:                     int64(1351727),
			ProfileCode:              &profileCode,
			GCCProfileCode:           &profileCode,
			Name:                     &name,
			Platform:                 "INSTAGRAM",
			Handle:                   &handle,
			Username:                 &username,
			PlatformThumbnail:        &thumbnail,
			EngagementRatePercenatge: 0.5445210507630777,
			Followers:                &followers,
			Following:                &following,
			AvgLikes:                 &avgLikes,
			AvgComments:              &avgComments,
			Uploads:                  &uploads,
			AvgReelsPlay30d:          &avgReelsPlay30d,
			AvgReach:                 &avgReach,
			StoryReach:               &storyReach,
			ImageReach:               &imageReach,
			ReelsReach:               &reelsReach,
			IsVerified:               &isVerified,
			Source:                   "GCC",
			AccountType:              &accountType,
			AudienceCityAvailable:    &audienceCityAvailable,
			AudienceGenderAvailable:  &audienceGenderAvailable,
		},
	}
}

// -------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------
func mockYaBeatResponse() *beatservice.CompleteProfileResponse {
	return &beatservice.CompleteProfileResponse{
		CompleteProfile: beatservice.CompleteProfile{
			Platform:           "youtube",
			YoutubeChannelId:   "UCq-Fj5jknLsUf-MWSy4_brA",
			YoutubeTitle:       "T-Series",
			YoutubeSubscribers: 249000000,
			YoutubeUploads:     19551,
			YoutubeViews:       233381716303,
			YoutubeThumbnail:   "https://yt3.ggpht.com/y1F4EOGuP19nZcBlzcyCtnHiYhkAOPQiRxwKeaGrOjXarUZZjcx_heiDiC06_Qj6ERea_qWK9A=s88-c-k-c0x00ffffff-no-rj",
		},
	}
}

func prepareYaCampaignProfileForTest() *coredomain.CampaignProfileEntry {
	profileCode := "IA_1351868"
	gccProfileCode := "YA_228252"
	name := "T-Series"
	handle := "UCq-Fj5jknLsUf-MWSy4_brA"
	username := "@tseries"
	thumbnail := "https://yt3.ggpht.com/y1F4EOGuP19nZcBlzcyCtnHiYhkAOPQiRxwKeaGrOjXarUZZjcx_heiDiC06_Qj6ERea_qWK9A=s512-c-k-c0x00ffffff-no-rj"
	followers := int64(246000000)
	uploads := int64(20117)
	views := int64(229163934559)
	avgReach := float64(0)
	audienceCityAvailable := false
	audienceGenderAvailable := true
	return &coredomain.CampaignProfileEntry{
		SocialDetails: &coredomain.SocialAccount{
			Code:                     int64(228252),
			ProfileCode:              &profileCode,
			GCCProfileCode:           &gccProfileCode,
			Name:                     &name,
			Platform:                 "YOUTUBE",
			Handle:                   &handle,
			Username:                 &username,
			PlatformThumbnail:        &thumbnail,
			EngagementRatePercenatge: 0.5473748164179614,
			Followers:                &followers,
			Uploads:                  &uploads,
			AvgReach:                 &avgReach,
			Source:                   "GCC",
			AudienceCityAvailable:    &audienceCityAvailable,
			AudienceGenderAvailable:  &audienceGenderAvailable,
			Views:                    &views,
		},
	}
}

func mockYaCpEntryToCheckBeatResponse() *coredomain.CampaignProfileEntry {
	profileCode := "IA_1351868"
	gccProfileCode := "YA_228252"
	name := "T-Series"
	handle := "UCq-Fj5jknLsUf-MWSy4_brA"
	username := "@tseries"
	thumbnail := "https://yt3.ggpht.com/y1F4EOGuP19nZcBlzcyCtnHiYhkAOPQiRxwKeaGrOjXarUZZjcx_heiDiC06_Qj6ERea_qWK9A=s88-c-k-c0x00ffffff-no-rj"
	followers := int64(249000000)
	uploads := int64(19551)
	views := int64(233381716303)
	avgReach := float64(0)
	audienceCityAvailable := false
	audienceGenderAvailable := true
	return &coredomain.CampaignProfileEntry{
		SocialDetails: &coredomain.SocialAccount{
			Code:                     int64(228252),
			ProfileCode:              &profileCode,
			GCCProfileCode:           &gccProfileCode,
			Name:                     &name,
			Platform:                 "YOUTUBE",
			Handle:                   &handle,
			Username:                 &username,
			PlatformThumbnail:        &thumbnail,
			EngagementRatePercenatge: 0.5473748164179614,
			Followers:                &followers,
			Uploads:                  &uploads,
			AvgReach:                 &avgReach,
			Source:                   "GCC",
			AudienceCityAvailable:    &audienceCityAvailable,
			AudienceGenderAvailable:  &audienceGenderAvailable,
			Views:                    &views,
		},
	}
}

func prepareYaProfileForTest() *coredomain.Profile {
	code := "IA_1351868"
	gccCode := "YA_228252"
	platformCode := "228252"
	name := "T-Series"
	handle := "UCq-Fj5jknLsUf-MWSy4_brA"
	username := "@tseries"
	description := "\"Music can change the world\". T-Series is India's largest Music Label \\u0026 Movie Studio, believes in bringing world close together through its music.\\nT-Series is associated with music industry from past three decades, having ample catalogue of music comprising plenty of languages that covers the length \\u0026 breadth of India. We believe after silence, nearest to expressing the inexpressible is Music. So, all the music lovers who believe in magic of music come join us and live the magic of music with T-Series."
	profileType := "Mega"
	thumbnail := "https://yt3.ggpht.com/y1F4EOGuP19nZcBlzcyCtnHiYhkAOPQiRxwKeaGrOjXarUZZjcx_heiDiC06_Qj6ERea_qWK9A=s512-c-k-c0x00ffffff-no-rj"
	followers := int64(250000000)
	uploads := int64(19551)
	avgViews := float64(11391556.124620967)
	videoReach := int64(11391556)
	views := int64(233734146988)
	return &coredomain.Profile{
		Code:         &code,
		Locked:       false,
		GccCode:      &gccCode,
		PlatformCode: platformCode,
		Name:         &name,
		Handle:       &handle,
		Username:     &username,
		Description:  &description,
		Platform:     "YOUTUBE",
		ProfileType:  &profileType,
		Thumbnail:    &thumbnail,
		Metrics: coredomain.ProfileMetrics{
			Followers:                &followers,
			EngagementRatePercenatge: 0.9881531070037496,
			AvgViews:                 &avgViews,
			VideoReach:               &videoReach,
			Uploads:                  &uploads,
			Views:                    &views,
		},
	}
}

func mockYAEntityToCheckBeatResponse() discoverydao.YoutubeAccountEntity {
	title := "T-Series"
	channelId := "UCq-Fj5jknLsUf-MWSy4_brA"
	profileType := "Mega"
	thumbnail := "https://yt3.ggpht.com/y1F4EOGuP19nZcBlzcyCtnHiYhkAOPQiRxwKeaGrOjXarUZZjcx_heiDiC06_Qj6ERea_qWK9A=s88-c-k-c0x00ffffff-no-rj"
	followers := int64(249000000)
	uploads := int64(19551)
	views := int64(233381716303)
	profileId := "IA_1351868"
	gccProfileId := "YA_228252"
	username := "@tseries"
	description := "\"Music can change the world\". T-Series is India's largest Music Label \\u0026 Movie Studio, believes in bringing world close together through its music.\\nT-Series is associated with music industry from past three decades, having ample catalogue of music comprising plenty of languages that covers the length \\u0026 breadth of India. We believe after silence, nearest to expressing the inexpressible is Music. So, all the music lovers who believe in magic of music come join us and live the magic of music with T-Series."
	videoReach := int64(11391556)
	avgViews := float64(11391556.124620967)
	return discoverydao.YoutubeAccountEntity{
		ID:           228252,
		ProfileId:    &profileId,
		GccProfileId: &gccProfileId,
		Username:     &username,
		Description:  &description,
		VideoReach:   &videoReach,
		AvgViews:     &avgViews,
		ChannelId:    &channelId,
		Title:        &title,
		Thumbnail:    &thumbnail,
		UploadsCount: &uploads,
		Followers:    &followers,
		ViewsCount:   &views,
		Label:        &profileType,
	}
}
