package points

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"receipt-processor/models"
)

func CalculatePoints(receipt models.Receipt) int {
    points := 0

    // 1. Points for alphanumeric characters in the retailer name
    re := regexp.MustCompile(`[a-zA-Z0-9]`)
    points += len(re.FindAllString(receipt.Retailer, -1))

    // 2. Points if the total is a round dollar amount
    total, _ := strconv.ParseFloat(receipt.Total, 64)
    if total == math.Floor(total) {
        points += 50
    }

    // 3. Points if the total is a multiple of 0.25
    if math.Mod(total, 0.25) == 0 {
        points += 25
    }

    // 4. Points for every two items
    points += (len(receipt.Items) / 2) * 5

    // 5. Points for items with description lengths that are multiples of 3
    for _, item := range receipt.Items {
        if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
            price, _ := strconv.ParseFloat(item.Price, 64)
            points += int(math.Ceil(price * 0.2))
        }
    }

    // 6. Points if the purchase day is odd
    day, _ := strconv.Atoi(receipt.PurchaseDate[8:10])
    if day%2 != 0 {
        points += 6
    }

    // 7. Points if the time of purchase is after 2:00pm and before 4:00pm
    purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
    if purchaseTime.Hour() == 14 || (purchaseTime.Hour() == 15 && purchaseTime.Minute() < 60) {
        points += 10
    }

    return points
}
