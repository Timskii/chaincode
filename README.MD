## Репозитарий для Chaincode


## Тесты
для запуска необходимы исходники hyperledger fabric, которые надо указать по пути `$GOPATH/src/github.com/hyperledger/fabric`

запуск тестов `test.sh` 

если при запуске тестов возникла ошибка:

`# github.com/hyperledger/fabric/vendor/github.com/miekg/pkcs11
../../../../vendor/github.com/miekg/pkcs11/pkcs11.go:26:18: fatal error: ltdl.h: No such file or directory
compilation terminated.`

надо скачать библиотекy:
`sudo apt install libltdl-dev`

---

в директориях chaincode имеется скрипт `./install_cc.sh` , для автоматической инсталяции их и ведения версии

так же имеется скрипт в корне install_all_cc.sh, для инсталяции всех Chaincode