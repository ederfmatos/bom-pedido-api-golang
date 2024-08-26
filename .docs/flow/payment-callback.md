# Fluxo de execução do recebiento de callback de pagamento

```mermaid
 flowchart TB
    A1[Mensagem recebida na fila de callback] --> A2{O que aconteceu?}
    A2 -- Pagamento Criado --> A3[Retorna]
    A2 -- Pagamento Alterado --> B1[Busca o pedido no banco de dados]
    B1 --> B2{Pedido existe?}
    B2 -- Não --> B3[Retorna]
    B2 -- Sim --> C1{Qual o estado do Pedido?}
    C1 -- Cancelado --> C2[Retorna]
    C1 -- Finalizado --> C3[Retorna]
    C1 -- Rejeitado --> C4[Retorna]
    C1 -- Aprovado --> C5[Retorna]
    C1 -- Em andamento --> C6[Retorna]
    C1 -- Aguardando retirada --> C7[Retorna]
    C1 -- Aguardando entrega --> C8[Retorna]
    C1 -- Em rota de entrega --> C9[Retorna]
    C1 -- Aguardando Pagamento --> D1[Busca a transação com status Pendente]
    D1 --> D2{Transação existe?}
    D2 -- Não --> D3[Retorna]
    D2 -- Sim --> E1[Consulta o Gateway de pagamento]
    
    E1 --> F1{Qual o status da transação?}
    
    F1 -- Pendente --> F2[Retorna]
    F1 -- Cancelada --> F3[Retorna]
    F1 -- Estornada --> F4[Retorna]
    F1 -- Paga --> G1[Atualiza pedido para AguardandoAprovação]
    G1 --> G2[Atualiza transação de Pendente para Paga]
```