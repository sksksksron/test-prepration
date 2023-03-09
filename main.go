package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	playerScoreRanks := PlayerScoreRank(os.Args[1])
	for _, data := range playerScoreRanks {
		fmt.Println(data)
	}
}

func PlayerScoreRank(fileName string) []string {
	rows, err := readInputCSVFile(fileName)
	if err != nil {
		panic(err)
	}

	playerMeanScores := calcMeanScoreByPlayer(rows)
	playerMeanScores = addRankPlayerMeanScore(playerMeanScores)

	outputCSV := outputCSVPlayerScore(playerMeanScores)

	return outputCSV
}

type ReadCSVRow struct {
	CreateTimeStamp time.Time
	PlayerID        string
	Score           int
}

func readInputCSVFile(fileName string) ([]ReadCSVRow, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	readCSVRows := []ReadCSVRow{}
	for i, row := range rows {
		// headerはスキップ
		// headerが意図したフォーマットかどうかチェックする必要があるので将来的に対応する。
		if i == 0 {
			continue
		}
		timeStampText := row[0]
		timeStamp, err := time.Parse("2006/01/02 15:04", timeStampText)
		if err != nil {
			return nil, err
		}
		playerID := row[1]
		scoreText := row[2]
		score, err := strconv.Atoi(scoreText)
		if err != nil {
			return nil, err
		}
		readCSVRows = append(readCSVRows, ReadCSVRow{
			CreateTimeStamp: timeStamp,
			PlayerID:        playerID,
			Score:           score,
		})

	}
	return readCSVRows, nil
}

type OutputPlayerMeanScore struct {
	Rank      int
	PlayerID  string
	MeanScore int
}

func calcMeanScoreByPlayer(rows []ReadCSVRow) []OutputPlayerMeanScore {
	playerMeanScores := []OutputPlayerMeanScore{}

	scoresByPlayer := make(map[string][]int)
	for _, row := range rows {
		playerID := row.PlayerID
		score := row.Score
		scoresByPlayer[playerID] = append(scoresByPlayer[playerID], score)
	}

	for player, scores := range scoresByPlayer {
		var scoreSum int
		for _, score := range scores {
			scoreSum = scoreSum + score
		}
		meanScore := scoreSum / len(scores)
		playerMeanScores = append(playerMeanScores, OutputPlayerMeanScore{
			PlayerID:  player,
			MeanScore: meanScore,
		})
	}
	return playerMeanScores
}

func addRankPlayerMeanScore(scores []OutputPlayerMeanScore) []OutputPlayerMeanScore {
	// 平均点をもとに降順に並び替え
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].MeanScore > scores[j].MeanScore
	})

	for i, score := range scores {
		// 同一平均スコアの場合は同一ランク
		if i > 0 && scores[i-1].MeanScore == score.MeanScore {
			scores[i].Rank = scores[i-1].Rank
			continue
		}
		scores[i].Rank = i + 1
	}
	return scores
}

func outputCSVPlayerScore(playerMeanScores []OutputPlayerMeanScore) []string {
	outputCSV := []string{}
	outputCSV = append(outputCSV, "rank,player_id,mean_score")
	for _, pms := range playerMeanScores {
		rank := strconv.Itoa(pms.Rank)
		playerID := pms.PlayerID
		meanScore := strconv.Itoa(pms.MeanScore)

		row := []string{rank, playerID, meanScore}
		outputCSV = append(outputCSV, strings.Join(row, ","))
	}
	return outputCSV
}
