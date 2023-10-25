package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const serverSalt = "this_is_a_salt"

var clientAnswerTimeout int = 5

var wordOfWisdom []string = []string{
	"A Word of Wisdom, for the benefit of the council of high priests, assembled in Kirtland, and the church, and also the saints in Zion—",
	"To be sent greeting; not by commandment or constraint, but by revelation and the word of wisdom, showing forth the order and will of God in the temporal salvation of all saints in the last days—",
	"Given for a principle with promise, adapted to the capacity of the weak and the weakest of all saints, who are or can be called saints.",
	"Behold, verily, thus saith the Lord unto you: In consequence of evils and designs which do and will exist in the hearts of conspiring men in the last days, I have warned you, and forewarn you, by giving unto you this word of wisdom by revelation—",
	"That inasmuch as any man drinketh wine or strong drink among you, behold it is not good, neither meet in the sight of your Father, only in assembling yourselves together to offer up your sacraments before him.",
	"And, behold, this should be wine, yea, pure wine of the grape of the vine, of your own make.",
	"And, again, strong drinks are not for the belly, but for the washing of your bodies.",
	"And again, tobacco is not for the body, neither for the belly, and is not good for man, but is an herb for bruises and all sick cattle, to be used with judgment and skill.",
	"And again, hot drinks are not for the body or belly.",
	"And again, verily I say unto you, all wholesome herbs God hath ordained for the constitution, nature, and use of man—",
	"Every herb in the season thereof, and every fruit in the season thereof; all these to be used with prudence and thanksgiving.",
	"Yea, flesh also of beasts and of the fowls of the air, I, the Lord, have ordained for the use of man with thanksgiving; nevertheless they are to be used sparingly;",
	"And it is pleasing unto me that they should not be used, only in times of winter, or of cold, or famine.",
	"All grain is ordained for the use of man and of beasts, to be the staff of life, not only for man but for the beasts of the field, and the fowls of heaven, and all wild animals that run or creep on the earth;",
	"And these hath God made for the use of man only in times of famine and excess of hunger.",
	"All grain is good for the food of man; as also the fruit of the vine; that which yieldeth fruit, whether in the ground or above the ground—",
	"Nevertheless, wheat for man, and corn for the ox, and oats for the horse, and rye for the fowls and for swine, and for all beasts of the field, and barley for all useful animals, and for mild drinks, as also other grain.",
	"And all saints who remember to keep and do these sayings, walking in obedience to the commandments, shall receive health in their navel and marrow to their bones;",
	"And shall find wisdom and great treasures of knowledge, even hidden treasures;",
	"And shall run and not be weary, and shall walk and not faint.",
}

func main() {
	if os.Getenv("CLIENT_ANSWER_TIMEOUT") != "" {
		timeout, err := strconv.Atoi(os.Getenv("CLIENT_ANSWER_TIMEOUT"))
		if err != nil {
			panic(err)
		}
		clientAnswerTimeout = timeout
	}
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func generatePrefix(host string) string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(randomBytes, rand.Uint64())
	return fmt.Sprintf("%d:%s:%s:%s", timestamp, host, serverSalt, base64.StdEncoding.EncodeToString(randomBytes))
}

func handleConnection(conn net.Conn) {
	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return
	}
	defer func() {
		fmt.Printf("Closing connection with %s\n", host)
		conn.Close()
	}()

	prefix := generatePrefix(host)
	conn.Write([]byte(prefix + "\n"))

	conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(clientAnswerTimeout)))

	reader := bufio.NewReader(conn)
	answer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	answer = strings.TrimRight(answer, "\r\n")
	fmt.Printf("answer: %s\n", answer)

	sum := sha256.Sum256([]byte(answer))
	fmt.Printf("hash: %x\n", sum)

	if !(sum[0] == 0 && sum[1] == 0 && (sum[2]&0xf0) == 0) {
		conn.Write([]byte("Hash doesn't begin with 5 zeros\n"))
		return
	}

	if !strings.HasPrefix(answer, prefix) {
		conn.Write([]byte("Prefix doesn't match\n"))
		return
	}

	randomQuote := wordOfWisdom[rand.Intn(len(wordOfWisdom))]
	conn.Write([]byte(randomQuote + "\n"))
}
