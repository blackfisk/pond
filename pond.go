package main

import (
        "code.google.com/p/go.crypto/openpgp"
        "code.google.com/p/go.crypto/openpgp/armor"
        "bufio"
        "fmt"
        "log"
        "os"
)

func getKeyByEmail(keyring openpgp.EntityList, email string) *openpgp.Entity {
        for _, entity := range keyring {
                if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
                        fmt.Print("Enter Passphrase: ")
                        bio := bufio.NewReader(os.Stdin)
                        line, _, _ := bio.ReadLine()

                        err := entity.PrivateKey.Decrypt(line)
                        if err != nil {
                                log.Println("asdasd")
                        }

                        for _, subkey := range entity.Subkeys {
                                if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
                                        err := subkey.PrivateKey.Decrypt(line)

                                        if err != nil {
                                                log.Println("sub")
                                        }
                                }
                        }
                }

                for _, ident := range entity.Identities {
                        log.Println(ident.UserId.Email)
                        if ident.UserId.Email == email {
                                return entity
                        }
                }
        }

        return nil
}

func main() {
        pubringFile, _ := os.Open("/Users/elcuervo/.gnupg/pubring.gpg")
        pubring, _ := openpgp.ReadKeyRing(pubringFile)
        privringFile, _ := os.Open("/Users/elcuervo/.gnupg/secring.gpg")
        privring, _ := openpgp.ReadKeyRing(privringFile)

        myPrivateKey := getKeyByEmail(privring, "yo@brunoaguirre.com")
        theirPublicKey := getKeyByEmail(pubring, "yo@brunoaguirre.com")

        w, _ := armor.Encode(os.Stdout, "PGP MESSAGE", nil)
        plaintext, _ := openpgp.Encrypt(w, []*openpgp.Entity{theirPublicKey}, myPrivateKey, nil, nil)
        fmt.Println("Write the message:")
        bio := bufio.NewReader(os.Stdin)
        message, _, _ := bio.ReadLine()

        fmt.Fprintf(plaintext, string(message))
        plaintext.Close()
        w.Close()
        fmt.Printf("\n")
}
