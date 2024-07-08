package fake

import (
	"math/rand"
	"strings"
	"time"

	"github.com/monopole/snips/internal/types"
)

func MakeSliceOfFakeUserData() []*types.MyUser {
	return []*types.MyUser{makeFakeUserData()}
}

func makeFakeUserData() *types.MyUser {
	repos := makeRandomRepoIds(6, 12)
	return &types.MyUser{
		Name:    "Wile E Coyote",
		Company: "Acme Corp.",
		Login:   "wcoyote",
		Email:   "wcoyote@acme.com",
		GhOrgs: []types.MyGhOrg{
			{Name: "tooFast", Login: "roadrunner"},
			{Name: "tooHigh", Login: "cliff"},
			{Name: "tooHeavy", Login: "anvil"},
		},
		IssuesCreated: makeIssueSet(
			makeRandomRepoIdGenerator(repos), 3+rand.Intn(5)),
		IssuesClosed: makeIssueSet(
			makeRandomRepoIdGenerator(repos), 3+rand.Intn(2)),
		IssuesCommented: makeIssueSet(
			makeRandomRepoIdGenerator(repos), 3+rand.Intn(8)),
		PrsReviewed: makeIssueSet(
			makeRandomRepoIdGenerator(repos), 3+rand.Intn(8)),
		Commits: makeCommitMap(
			makeRandomRepoIdGenerator(repos), 3+rand.Intn(8)),
	}
}

func makeCommitMap(repoIdGen *randomRepoIdGenerator, count int) map[types.RepoId][]*types.MyCommit {
	result := make(map[types.RepoId][]*types.MyCommit)
	for i := 0; i < count; i++ {
		commits := makeSliceOfCommits(3 + rand.Intn(15))
		repoId := *repoIdGen.get()
		for j := range commits {
			commits[j].RepoId = repoId
		}
		result[repoId] = commits
	}
	return result
}

func makeSliceOfCommits(count int) []*types.MyCommit {
	result := make([]*types.MyCommit, count)
	for i := 0; i < count; i++ {
		result[i] = makeRandomCommit()
	}
	return result
}

func makeRandomCommit() *types.MyCommit {
	issue := makeRandomIssue()
	return &types.MyCommit{
		// RepoId:           types.RepoId{},
		Sha:              "a123das3d1",
		Url:              "http://www.commit.com",
		MessageFirstLine: randLorem.getOkToReUse(),
		Committed:        time.Now().Add(randNegativeDay()),
		Author:           randUserName.getOkToReUse(),
		Pr:               &issue,
	}
}

func randNegativeDay() time.Duration {
	const day = 24 * time.Hour
	dayCount := -(1 + rand.Intn(365))
	return time.Duration(dayCount) * day
}

func makeIssueSet(repoIdGen *randomRepoIdGenerator, count int) *types.IssueSet {
	result := types.IssueSet{
		Domain: "github.acme.com",
		Groups: make(map[types.RepoId][]types.MyIssue),
	}
	for i := 0; i < count; i++ {
		issues := makeSliceOfIssues(1 + rand.Intn(15))
		repoId := repoIdGen.get()
		for j := range issues {
			issues[j].RepoId = *repoId
		}
		result.Groups[*repoId] = issues
	}
	return &result
}

func makeSliceOfIssues(count int) []types.MyIssue {
	result := make([]types.MyIssue, count)
	for i := 0; i < count; i++ {
		result[i] = makeRandomIssue()
	}
	return result
}

func makeRandomIssue() types.MyIssue {
	return types.MyIssue{
		//	RepoId:  types.RepoId{},
		Number:  rand.Intn(10000),
		Title:   randLorem.getOkToReUse(),
		HtmlUrl: "http://www.example.com",
		Updated: time.Now().Add(randNegativeDay()),
	}
}

func makeRandomRepoIds(numOrgs, numReposInOrg int) []*types.RepoId {
	var result []*types.RepoId
	for i := 0; i < numOrgs; i++ {
		org := "org" + randFruit.getUnique()
		for j := 0; j < numReposInOrg; j++ {
			result = append(result, &types.RepoId{
				Org:  org,
				Name: "repo" + randElement.getUnique(),
			})
		}
	}
	return result
}

type randomRepoIdGenerator struct {
	index  int
	source []*types.RepoId
}

func makeRandomRepoIdGenerator(raw []*types.RepoId) *randomRepoIdGenerator {
	return &randomRepoIdGenerator{
		index:  -1,
		source: raw,
	}
}

func (rs *randomRepoIdGenerator) get() *types.RepoId {
	if rs.index+1 == len(rs.source) {
		panic("everything has been used")
	}
	rs.index++
	return rs.source[rs.index]
}

func shuffle(src []string) []string {
	dest := make([]string, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}

type randomStringGenerator struct {
	index  int
	source []string
}

func makeRandomStringGenerator(raw []string) *randomStringGenerator {
	return &randomStringGenerator{
		index:  -1,
		source: shuffle(raw),
	}
}

func (rs *randomStringGenerator) getUnique() string {
	if rs.index+1 == len(rs.source) {
		panic("depleted " + strings.Join(rs.source, ","))
	}
	rs.index++
	return rs.source[rs.index]
}

func (rs *randomStringGenerator) getOkToReUse() string {
	return rs.source[rand.Intn(len(rs.source))]
}

var (
	randLorem = makeRandomStringGenerator(strings.Split(`
Lorem ipsum dolor sit amet
consectetur adipiscing elit
sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua
Ut enim ad minim veniam
quis nostrud exercitation ullamco
laboris nisi ut
aliquip ex ea commodo consequat
Duis aute irure dolor in reprehenderit
in voluptate velit esse cillum dolore eu
fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non
proident sunt in culpa qui officia
deserunt mollit anim id est laborum`[1:], "\n"))

	randUserName = makeRandomStringGenerator(strings.Split(`
Noah
Emma
Oliver
Charlotte
Elijah
Amelia
James
Ava
William
Sophia
Benjamin
Isabella
Lucas
Mia
Henry
Evelyn`[1:], "\n"))

	randFruit = makeRandomStringGenerator(strings.Split(`
Apple
Banana
Grape
Melon
Avocado
Carrot
Pepper
Pineapple
Cherry
Pear
Lemon
Lime
Kiwi
Grapefruit
Mango
Pomegranate
Tangerine`[1:], "\n"))

	randElement = makeRandomStringGenerator(strings.Split(`
Hydrogen
Helium
Lithium
Beryllium
Boron
Carbon
Nitrogen
Oxygen
Fluorine
Neon
Sodium
Magnesium
Aluminium
Silicon
Phosphorus
Sulfur
Chlorine
Argon
Potassium
Calcium
Scandium
Titanium
Vanadium
Chromium
Manganese
Iron
Cobalt
Nickel
Copper
Zinc
Gallium
Germanium
Arsenic
Selenium
Bromine
Krypton
Rubidium
Strontium
Yttrium
Zirconium
Niobium
Molybdenum
Technetium
Ruthenium
Rhodium
Palladium
Silver
Cadmium
Indium
Tin
Antimony
Tellurium
Iodine
Xenon
Cesium
Barium
Lanthanum
Cerium
Praseodymium
Neodymium
Promethium
Samarium
Europium
Gadolinium
Terbium
Dysprosium
Holmium
Erbium
Thulium
Ytterbium
Lutetium
Hafnium
Tantalum
Tungsten
Rhenium
Osmium
Iridium
Platinum
Gold
Mercury
Thallium
Lead
Bismuth
Polonium
Astatine
Radon
Francium
Radium
Actinium
Thorium
Protactinium
Uranium
Neptunium
Plutonium
Americium
Curium
Berkelium
Californium
Einsteinium
Fermium
Mendelevium
Nobelium
Lawrencium
Rutherfordium
Dubnium
Seaborgium
Bohrium
Hassium
Meitnerium
Darmstadtium
Roentgenium
Copernicium
Nihonium
Flerovium
Moscovium
Livermorium
Tennessine
Oganesson`[1:], "\n"))
)
