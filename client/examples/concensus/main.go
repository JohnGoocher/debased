package main

//used for consensus protocol
import (
	// TODO: import blockchain and related structs as bc
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	//"bytes"
	"os"
	"strings"
	"encoding/json"
)
