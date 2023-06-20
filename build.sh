#!/bin/bash

# Executa serial.sh
echo "Executando o primeiro script...Versão sequencial - somente uma thread."
chmod +x serial.sh
./serial.sh

# Executa concurrent.sh
echo "Executando o segundo script...Versão concorrente - várias threads."
chmod +x concurrent.sh
./concurrent.sh

echo "Finalizado!"
