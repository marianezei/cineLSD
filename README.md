# cineLSD


Objetivo

Neste laboratório, precisa-se escrever um programa que interage com serviço web. Este serviço mantém informações sobre atores e filmes nos quais estes atores atuaram. O objetivo principal do seu programa deve ser montar um ranking de atores. O valor do score de um ator é calculado com base no rating dos filmes que este autor participou.

Como sugestão, para atingir esse objetivo, você deve considerar as seguintes etapas:
Coleta de dados de todos os atores (verifique um dump dos IDs do dataset de atores aqui);
Coleta de dados dos filmes associados com os atores;
Cálculo do score do ator (média aritmética dos ratings dos filmes que o ator participou);
Rankeamento do atores;
Saída do programa com os top-10 atores.

Você deve implementar duas versões do seu programa. Uma totalmente sequencial (com somente uma thread) e outra concorrente. O objetivo é que a versão concorrente seja mais eficiente (retorna os top-10 atores em menos tempo que a versão sequencial).
