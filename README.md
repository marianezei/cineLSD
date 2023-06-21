# cineLSD


## Objetivo

Neste laboratório, precisa-se escrever um programa que interage com serviço web. Este serviço mantém informações sobre atores e filmes nos quais estes atores atuaram. O objetivo principal do seu programa é montar um ranking de atores. O valor do score de um ator é calculado com base no rating dos filmes que este ator participou.

Existem duas versões do programa. Uma totalmente sequencial (com somente uma thread) e outra concorrente. O objetivo é que a versão concorrente seja mais eficiente (retorna os top-10 atores em menos tempo que a versão sequencial).

## Estratégia (Versão Concorrente)

A versão concorrente foi criada a partir da versão sequencial do nosso código, com uma mudança: o código foi dividido em funções, facilitando a implementação do paralelismo.

Na função **calculateActorScores** nós definimos:

1. Uma constante que indica o tamanho (16) da nossa workerPool
2. Uma variável do tipo sync.waitGroup
3. Uma variável do tipo sync.Mutex

1. Na nossa primeira implementação da versão concorrente sequer existia a ideia de workerPool. Uma [goroutine](https://go.dev/tour/concurrency/1) era criada para cada ator do arquivo [actors.txt](./actors.txt), o que se mostrou problemático, visto que o tamanho do arquivo original era muito grande e o grande número de threads poderia prejudicar a performance do sistema. Então, passamos a fazer uso de uma workerPool e definimos o seu valor por meio de testes manuais, 16 goroutines/threads se comportou como um número ideal, acima desse valor pouco ganhamos em performance. 

2. O [waitGroup](https://go.dev/src/sync/waitgroup.go) foi utilizado para garantir que TODAS as goroutines terminaram de executar antes de retornar o resultado.

3. Por fim, usamos [Mutex](https://go.dev/tour/concurrency/9) para garantir a segurança das variáveis compartilhadas, evitando condições de corrida. 

Abaixo, a parte concorrente do código comentada:

```go
func calculateActorScores(actorIDs []string) (map[string]float32, error) {
	actorScores := make(map[string]float32) //Variável compartilhada
	actorMovieCount := make(map[string]int) //Variável compartilhada

	const numWorkers = 16
	var wg sync.WaitGroup
	var mutex sync.Mutex

	actorCh := make(chan string) //Canal usado para compartilhar os ids com as threads

	worker := func() { //Define a função do Worker
		defer wg.Done() //Decrementa o número de workers sempre que a goroutine encerra

		for actorID := range actorCh {
			actor, err := getActor(actorID)
			if err != nil {
				log.Printf("Falha ao obter informações do ator %s: %v", actorID, err)
				continue
			}

			var totalScore float32
			for _, movieID := range actor.Movies {
				movie, err := getMovie(strings.ReplaceAll(movieID, `"`, ""))
				if err != nil {
					log.Printf("Falha ao obter informações do filme %s: %v", movieID, err)
					continue
				}

				totalScore += movie.Score //Variável local

				mutex.Lock()
				actorMovieCount[actor.Name]++ //Variável compartilhada protegida
				mutex.Unlock()
			}

			mutex.Lock()
			if actorMovieCount[actor.Name] > 0 { // Variável compartilhada protegida
				actorScores[actor.Name] = totalScore / float32(actorMovieCount[actor.Name])
			} else {
				actorScores[actor.Name] = 0
			}
			mutex.Unlock()
		}
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)   //Adiciona a goroutine ao grupo de workers
		go worker() //Inicia a goroutine (worker)
	}

	for _, actorID := range actorIDs {
		actorCh <- actorID
	}
	close(actorCh)

	wg.Wait() //Espera por todos os workers do grupo finalizarem

	return actorScores, nil
}
```


