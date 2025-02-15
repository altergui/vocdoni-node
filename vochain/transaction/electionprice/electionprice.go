// Package electionprice provides a mechanism for calculating the price of an election based on its characteristics.
//
// The formula used to calculate the price for creating an election on the Vocdoni blockchain is designed to take into
// account various factors that impact the cost and complexity of conducting an election. The price is determined by
// combining several components, each reflecting a specific aspect of the election process.
//
// 1. Base Price: This is a fixed cost that serves as a starting point for the price calculation. It represents the
// minimal price for creating an election, regardless of its size or duration.
//
// 2. Size Price: As the number of voters (maxCensusSize) in an election increases, the resources required to manage
// the election also grow. To account for this, the size price component is directly proportional to the maximum number
// of votes allowed in the election. Additionally, it takes into consideration the blockchain's maximum capacity
// (capacity) and the maximum capacity the blockchain administrators can set (maxCapacity). This ensures that the price
// is adjusted based on the current capacity of the blockchain.
//
// 3. Duration Price: The length of the election (electionDuration) also affects the price, as longer elections occupy
// more resources over time. The duration price component is directly proportional to the election duration and
// inversely proportional to the maximum number of votes. This means that if the election lasts longer, the price
// increases, and if there are more votes in a shorter time, the price also increases to reflect the higher demand for
// resources.
//
// 4. Encrypted Votes: If an election requires encryption for maintaining secrecy until the end (encryptedVotes), it
// demands additional resources and computational effort. Therefore, the encrypted price component is added to the total
// price when this feature is enabled.
//
// 5. Anonymous Votes: Similarly, if an election must be anonymous (anonymousVotes), it requires additional measures to
// ensure voter privacy. As a result, the anonymous price component is added to the total price when this option is
// chosen.
//
// 6. Overwrite Price: Allowing voters to overwrite their votes (maxVoteOverwrite) can increase the complexity of
// managing the election, as it requires additional resources to handle vote updates. The overwrite price component
// accounts for this by being proportional to the maximum number of vote overwrites and the maximum number of votes
// allowed in the election. It also takes into account the blockchain's capacity to ensure the price reflects the
// current resource constraints.
//
// The constant factors in the price formula play a crucial role in determining the price of an election based on its
// characteristics. Each factor is associated with a specific component of the price formula and helps to weigh the
// importance of that component in the final price calculation. The rationale behind these constant factors is to
// provide a flexible mechanism to adjust the pricing model based on the system's needs and requirements.
//
// k1 (Size price factor): This constant factor affects the size price component of the formula. By adjusting k1,
// you can control the impact of the maximum number of votes (maxCensusSize) on the overall price. A higher k1 value
// would make the price increase more rapidly as the election size grows, while a lower k1 value would make the price
// less sensitive to the election size. The rationale behind k1 is to ensure that the pricing model can be adapted to
// accommodate different election sizes while considering the resource requirements.
//
// k2 (Duration price factor): This constant factor influences the duration price component of the formula. By
// adjusting k2, you can control how the duration of the election (electionDuration) affects the price. A higher k2
// value would make the price increase more quickly as the election duration extends, while a lower k2 value would make
// the price less sensitive to the election duration. The rationale behind k2 is to reflect the resource consumption
// over time and ensure that longer elections are priced accordingly.
//
// k3 (Encrypted price factor): This constant factor affects the encrypted price component of the formula. By adjusting
// k3, you can control the additional cost associated with encrypted elections (encryptedVotes). A higher k3 value would
// make the price increase more significantly for elections that require encryption, while a lower k3 value would make
// the price less sensitive to the encryption requirement. The rationale behind k3 is to account for the extra
// computational effort and resources needed to ensure secrecy in encrypted elections.
//
// k4 (Anonymous price factor): This constant factor influences the anonymous price component of the formula. By
// adjusting k4, you can control the additional cost associated with anonymous elections (anonymousVotes). A higher k4
// value would make the price increase more significantly for elections that require anonymity, while a lower k4 value
// would make the price less sensitive to the anonymity requirement. The rationale behind k4 is to account for the extra
// measures and resources needed to ensure voter privacy in anonymous elections.
//
// k5 (Overwrite price factor): This constant factor affects the overwrite price component of the formula. By adjusting
// k5, you can control the additional cost associated with allowing vote overwrites (maxVoteOverwrite). A higher k5
// value would make the price increase more significantly for elections that permit vote overwrites, while a lower k5
// value would make the price less sensitive to the overwrite allowance. The rationale behind k5 is to account for the
// increased complexity and resources needed to manage vote overwrites in the election process.
//
// k6 (Non-linear growth factor): This constant factor determines the rate of price growth for elections with a maximum
// number of votes (maxCensusSize) exceeding the k7 threshold. By adjusting k6, you can control the non-linear growth
// rate of the price for larger elections. A higher k6 value would result in a more rapid increase in the price as the
// election size grows beyond the k7 threshold, while a lower k6 value would result in a slower increase in the price
// for larger elections. The rationale behind k6 is to provide a mechanism for controlling the pricing model's
// sensitivity to large elections. This factor ensures that the price accurately reflects the increased complexity,
// resource consumption, and management effort associated with larger elections, while maintaining a more affordable
// price for smaller elections. By fine-tuning k6, the pricing model can be adapted to balance accessibility for smaller
// elections with the need to cover costs and resource requirements for larger elections.
//
// k7 (Size non-linear trigger): This constant factor represents a threshold value for the maximum number of
// votes (maxCensusSize) in an election. When the election size exceeds k7, the price growth becomes non-linear,
// increasing more rapidly beyond this point. The rationale behind k7 is to create a pricing model that accommodates
// a "freemium" approach, where smaller elections (under the k7 threshold) are priced affordably, while larger elections
// are priced more significantly due to their increased resource requirements and complexity. By adjusting k7, you can
// control the point at which the price transition from linear to non-linear growth occurs. A higher k7 value would
// allow for more affordable pricing for a larger range of election sizes, while a lower k7 value would result in more
// rapid price increases for smaller election sizes. This flexibility enables the pricing model to be tailored to the
// specific needs and goals of the Vocdoni blockchain, ensuring that small elections remain accessible and affordable,
// while larger elections are priced to reflect their higher resource demands.
package electionprice

import (
	"math"
)

// ElectionParameters is a struct to group the input parameters for CalculatePrice method.
type ElectionParameters struct {
	MaxCensusSize    int
	ElectionDuration int
	EncryptedVotes   bool
	AnonymousVotes   bool
	MaxVoteOverwrite int
}

// Calculator is a struct that stores the constant factors and basePrice
// required for calculating the price of an election.
type Calculator struct {
	basePrice uint32
	capacity  int
	factors   Factors
}

// Factors is a struct that stores the constant factors required for calculating
// the price of an election.
type Factors struct {
	k1 float64
	k2 float64
	k3 float64
	k4 float64
	k5 float64
	k6 float64
	k7 int
}

// DefaultElectionPriceFactors is the default set of constant factors used for calculating the price.
var DefaultElectionPriceFactors = Factors{
	k1: 0.0008,
	k2: 0.002,
	k3: 0.005,
	k4: 10,
	k5: 3,
	k6: 0.001,
	k7: 500,
}

// NewElectionPriceCalculator creates a new PriceCalculator with the given constant factors,
// basePrice, and maxCapacity.
func NewElectionPriceCalculator(basePrice uint32, capacity int,
	factors Factors) *Calculator {
	return &Calculator{
		basePrice: basePrice,
		capacity:  capacity,
		factors:   factors,
	}
}

// Price computes the price of an election given the parameters.
// The price is calculated using the following formula:
// price = basePrice + sizePrice + durationPrice + encryptedPrice + anonymousPrice + overwritePrice
//
// Parameters:
//   - MaxCensusSize: The maximum number of votes casted allowed for an election.
//   - ElectionDuration: The number of blocks the election can last. Currently the block time is 10 seconds.
//   - EncryptedVotes: A boolean flag that indicates if the election requires encryption keys for the
//     secret-until-the-end property.
//   - AnonymousVotes: A boolean flag that indicates if the election is anonymous or not.
//   - MaxVoteOverwrite: The number of overwrites a voter can execute after sending the first vote.
//
// The output is a uint64 value representing the price in a suitable unit.
func (p *Calculator) Price(params *ElectionParameters) uint64 {
	sizePrice := p.factors.k1 *
		float64(params.MaxCensusSize) *
		(1 - (float64(1 / p.capacity))) *
		(1 + p.factors.k6*math.Max(0, float64(params.MaxCensusSize-p.factors.k7)))

	durationPrice := p.factors.k2 *
		float64(params.ElectionDuration) *
		(1 + (float64(params.MaxCensusSize) / float64(p.capacity)))

	encryptedPrice := 0.0
	if params.EncryptedVotes {
		encryptedPrice = p.factors.k3 * float64(params.MaxCensusSize)
	}

	anonymousPrice := 0.0
	if params.AnonymousVotes {
		anonymousPrice = p.factors.k4
	}

	overwritePrice := p.factors.k5 *
		float64(params.MaxCensusSize) *
		float64(params.MaxVoteOverwrite) /
		float64(p.capacity)

	price := float64(p.basePrice) + sizePrice + durationPrice + encryptedPrice + anonymousPrice + overwritePrice
	return uint64(price)
}

// SetCapacity sets the current capacity of the blockchain.
func (p *Calculator) SetCapacity(capacity int) {
	p.capacity = capacity
}
