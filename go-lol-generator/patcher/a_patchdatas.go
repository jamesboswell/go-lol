package patcher

import (
	"go/types"

	"github.com/pkg/errors"
)

// default
var overrides = &Patches{
	Resources: make(map[string]ResPatch),
}

func init() {
	overrides.Add("lol-static-data", ResPatch{
		Operations: map[string]OpPatch{
			"/champion":            {Name: "Champions"},
			"/champion/{id}":       {Name: "Champion"},
			"/item":                {Name: "Items"},
			"/item/{id}":           {Name: "Item"},
			"/language-strings":    {Name: "LanguageStrings"},
			"/languages":           {Name: "Languages"},
			"/map":                 {Name: "Maps"},
			"/mastery":             {Name: "Masteries"},
			"/mastery/{id}":        {Name: "Mastery"},
			"/realm":               {Name: "Realm"},
			"/rune":                {Name: "Runes"},
			"/rune/{id}":           {Name: "Rune"},
			"/summoner-spell":      {Name: "SummonerSpells"},
			"/summoner-spell/{id}": {Name: "SummonerSpell"},
			"versions":             {Name: "Versions"},
		},
		Classes: map[string]ClassPatch{
			"ImageDto":             {Name: "Image"},
			"ChampionDto":          {Name: "Champion"},
			"ChampionListDto":      {Name: "Champions"},
			"MapDetailsDto":        {Name: "Map"},
			"MapDataDto":           {Name: "Maps"}, // rito pls..
			"ChampionSpellDto":     {Name: "ChampionSpell"},
			"SummonerSpellDto":     {Name: "SummonerSpell"},
			"SummonerSpellListDto": {Name: "SummonerSpells"},
			"ItemDto":              {Name: "Item"},
			"ItemListDto":          {Name: "Items"},
			"GoldDto":              {Name: "Gold"},
			"StatsDto":             {Name: "ChampionStats"},
			"GroupDto":             {Name: "ItemGroup"},
			"InfoDto":              {Name: "ChampionInfo"},
			"SkinDto":              {Name: "Skin"},
			"RecommendedDto":       {Name: "Recommended"},
			"BlockDto":             {Name: "RecommendedBlock"},
			"BlockItemDto":         {Name: "RecommendedItems"},
			"ItemTreeDto":          {Name: "ItemTree"},
			"MasteryDto":           {Name: "Mastery"},
			"MasteryListDto":       {Name: "Masteries"},
			"MasteryTreeItemDto":   {Name: "MasteryTreeItem"},
			"MasteryTreeDto":       {Name: "MasteryTree"},
			"MasteryTreeListDto":   {Name: "MasteryTrees"},
			"RuneDto":              {Name: "Rune"},
			"RuneListDto":          {Name: "Runes"},
			"MetaDataDto":          {Name: "RuneMetadata"},
			"PassiveDto":           {Name: "Passive"},
			"SpellVarsDto":         {Name: "SpellVars"},
			"BasicDataDto":         {Name: "BasicData"},
			"BasicDataStatsDto":    {Name: "BasicStats"},
			"LanguageStringsDto":   {Name: "LanguageStrings"},
			"LevelTipDto":          {Name: "LevelTip"},
			"RealmDto":             {Name: "Realm"},
		},
	})

	overrides.Add("champion", ResPatch{
		Operations: map[string]OpPatch{
			"/champion":      {Name: "ChampionStatuses"},
			"/champion/{id}": {Name: "ChampionStatus"},
		},
		Classes: map[string]ClassPatch{
			"ChampionDto":     {Name: "ChampionStatus"},
			"ChampionListDto": {Name: "ChampionStatuses"},
		},
	})

	overrides.Add("current-game", ResPatch{
		Operations: map[string]OpPatch{
			"/getSpectatorGameInfo/{platformId}/{summonerId}": {Name: "SpectatorGameInfo"},
		},
		Classes: map[string]ClassPatch{
			"BannedChampion":         {Name: "CurrentGameBannedChampion"},
			"Rune":                   {Name: "CurrentGameRune"},
			"Mastery":                {Name: "CurrentGameMastery"},
			"CurrentGameInfo":        {Name: "CurrentGameInfo"},
			"CurrentGameParticipant": {Name: "CurrentGameParticipant"},
			"Observer":               {Name: "CurrentGameObserver"},
		},
	})

	overrides.Add("featured-games", ResPatch{
		Operations: map[string]OpPatch{
			"/featured": {Name: "FeaturedGames"},
		},
		Classes: map[string]ClassPatch{
			"BannedChampion":   {Name: "FeaturedGameBannedChampion"},
			"Participant":      {Name: "FeaturedGameParticipant"},
			"Rune":             {Name: "FeaturedGameRune"},
			"Mastery":          {Name: "FeaturedGameMastery"},
			"Observer":         {Name: "FeaturedGameObserver"},
			"FeaturedGames":    {Name: "FeaturedGames"},
			"FeaturedGameInfo": {Name: "FeaturedGameInfo"},
		},
	})

	overrides.Add("game", ResPatch{
		Operations: map[string]OpPatch{
			"/game/by-summoner/{summonerId}/recent": {Name: "RecentGames"},
		},
		Classes: map[string]ClassPatch{
			"GameDto":        {Name: "Game"},
			"RecentGamesDto": {Name: "RecentGames"},
			"RawStatsDto":    {Name: "GamePlayerRawStats"},
			"PlayerDto":      {Name: "GamePlayer"},
		},
	})

	overrides.Add("league", ResPatch{
		Operations: map[string]OpPatch{
			"/league/by-summoner/{summonerIds}":       {Name: "LeaguesBySummonerID"},
			"/league/by-summoner/{summonerIds}/entry": {Name: "LeagueEntriesBySummonerID"},
			"/league/by-team/{teamIds}":               {Name: "LeaguesByTeamID"},
			"/league/by-team/{teamIds}/entry":         {Name: "LeagueEntriesByTeamID"},
			"/league/challenger":                      {Name: "Challenger"},
			"/league/master":                          {Name: "Master"},
		},
		Classes: map[string]ClassPatch{
			"MiniSeriesDto":  {Name: "MiniSeries"},
			"LeagueEntryDto": {Name: "LeagueEntry"},
			"LeagueDto":      {Name: "League"},
		},
	})

	overrides.Add("lol-status", ResPatch{
		Operations: map[string]OpPatch{
			"/shards":          {Name: "Shards"},
			"/shards/{region}": {Name: "ShardsInRegion"},
			"/shards/{shard}":  {Name: "Shard"},
		},
		Classes: map[string]ClassPatch{
			"Shard":       {Name: "Shard"},
			"ShardStatus": {Name: "ShardStatus"},
			"Service":     {Name: "Service"},
			"Message":     {Name: "StatusMessage"},
			"Translation": {Name: "StatusMessageTranslation"},
			"Incident":    {Name: "Incident"},
		},
	})

	overrides.Add("match", ResPatch{
		Operations: map[string]OpPatch{
			"/match/{matchId}":                          {Name: "Match"},
			"/match/by-tournament/{tournamentCode}/ids": {Name: "MatchesByTournement"},
			"/match/for-tournament/{matchId}":           {Name: "MatchForTournement"},
		},
		Classes: map[string]ClassPatch{
			"BannedChampion":          {Name: "BannedChampion"},
			"Timeline":                {Name: "Timeline"},
			"Frame":                   {Name: "Frame"},
			"Event":                   {Name: "Event"},
			"Position":                {Name: "Position"},
			"Team":                    {Name: "MatchTeam"},
			"Participant":             {Name: "Participant"},
			"ParticipantStats":        {Name: "ParticipantStats"},
			"ParticipantIdentity":     {Name: "ParticipantIdentity"},
			"ParticipantFrame":        {Name: "ParticipantFrame"},
			"ParticipantStatus":       {Name: "ParticipantStatus"},
			"ParticipantTimeline":     {Name: "ParticipantTimeline"},
			"ParticipantTimelineData": {Name: "ParticipantTimelineData"},
			"Player":                  {Name: "Player"},
			"Mastery":                 {Name: "UsedMastery"},
			"Rune":                    {Name: "UsedRune"},
			"MatchDetail":             {Name: "MatchDetail"},
		},
	})

	overrides.Add("matchlist", ResPatch{
		Operations: map[string]OpPatch{
			"/matchlist/by-summoner/{summonerId}": {Name: "MatchesBySummonerID"},
		},
		Classes: map[string]ClassPatch{
			"MatchList":      {Name: "Matches"},
			"MatchReference": {Name: "MatchRef"},
		},
	})

	overrides.Add("stats", ResPatch{
		Operations: map[string]OpPatch{
			"/stats/by-summoner/{summonerId}/ranked":  {Name: "RankedStats"},
			"/stats/by-summoner/{summonerId}/summary": {Name: "StatsSummary"},
		},
		Classes: map[string]ClassPatch{
			"RankedStatsDto":            {Name: "RankedStats"},
			"PlayerStatsSummaryDto":     {Name: "PlayerStatsSummary"},
			"PlayerStatsSummaryListDto": {Name: "PlayerStatsSummaries"},
			"AggregatedStatsDto":        {Name: "AggregatedStats"},
			"ChampionStatsDto":          {Name: "PlayerChampionStats"},
		},
	})

	overrides.Add("summoner", ResPatch{
		Operations: map[string]OpPatch{
			"/summoner/by-name/{summonerNames}": {Name: "SummonersByName"},
			"/summoner/{summonerIds}":           {Name: "Summoners", MapKey: types.Int64},
			"/summoner/{summonerIds}/masteries": {Name: "MasteryPages", MapKey: types.Int64},
			"/summoner/{summonerIds}/name":      {Name: "SummonerNames", MapKey: types.Int64},
			"/summoner/{summonerIds}/runes":     {Name: "RunePages", MapKey: types.Int64},
		},
		Classes: map[string]ClassPatch{
			"MasteryDto":      {Name: "EquippedMastery"},
			"MasteryPageDto":  {Name: "MasteryPage"},
			"MasteryPagesDto": {Name: "MasteryPages"},
			"RuneSlotDto":     {Name: "RuneSlot"},
			"RunePageDto":     {Name: "RunePage"},
			"RunePagesDto":    {Name: "RunePages"},
			"SummonerDto":     {Name: "Summoner"},
		},
	})

	overrides.Add("team", ResPatch{
		Operations: map[string]OpPatch{
			"/team/by-summoner/{summonerIds}": {Name: "TeamsBySummonerID", MapKey: types.Int64},
			"/team/{teamIds}":                 {Name: "Teams"},
		},
		Classes: map[string]ClassPatch{
			"TeamDto":                {Name: "Team"},
			"MatchHistorySummaryDto": {Name: "TeamMatchHistorySummary"},
			"TeamStatDetailDto":      {Name: "TeamStatDetails"},
			"TeamMemberInfoDto":      {Name: "TeamMemberInfo"},
			"RosterDto":              {Name: "TeamRoaster"},
		},
	})

	overrides.Add("championmastery", ResPatch{
		Operations: map[string]OpPatch{
			"/championmastery/location/{platformId}/player/{playerId}/champion/{championId}": {
				Name: "ChampionMastery",
			},
			"/championmastery/location/{platformId}/player/{playerId}/champions": {
				Name: "ChampionMasteries",
			},
			"/championmastery/location/{platformId}/player/{playerId}/score": {
				Name: "ChampionMasteryScore",
			},
			"/championmastery/location/{platformId}/player/{playerId}/topchampions": {
				Name: "TopChampions",
			},
		},
		Classes: map[string]ClassPatch{
			"ChampionMasteryDTO": {Name: "ChampionMastery"},
		},
	})

}

type Patches struct {
	Resources map[string]ResPatch
}

type ResPatch struct {
	// map[path suffix]Operation
	Operations map[string]OpPatch
	Classes    map[string]ClassPatch
}

// OpPatch represets a predeclared operation info.
type OpPatch struct {
	// Method name on client.
	//
	// Required
	Name string

	// Override map key in return value.
	// Patch will panic if original return value is not map.
	MapKey types.BasicKind
}

type ClassPatch struct {
	Name string
}

func (rp *ResPatch) Class(clsName string) (*ClassPatch, error) {
	co, ok := rp.Classes[clsName]
	if !ok {
		return nil, errors.Errorf("class override is required for %q", clsName)
	}
	return &co, nil
}

func (p *Patches) Add(id string, rp ResPatch) {
	if _, ok := p.Resources[id]; ok {
		panic(errors.Errorf("resource override already exists: %q", id))
	}
	// for path, oo := range rp.Operations {
	// 	ovs.mtdName.Take(id, path, &oo)
	// }

	// for origName, co := range rp.Classes {
	// 	ovs.clsName.Take(id, origName, &co)
	// }

	p.Resources[id] = rp
}
