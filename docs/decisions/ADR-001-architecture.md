# ADR-001: Arquitetura do Sistema de Controle de Chamados Internos

status: Proposto
data: 2026-09-21

# Context:

O cliente trabalha na área administrativa de uma empresa e entrou em contato com a gente para desenvolver uma aplicação web para gerenciar chamados de suporte. A aplicação deve permitir que os usuários registrem e acompanhem os seus chamados, enquanto os administradores e o time de suporte possam gerenciar e responder de forma eficiente e com uma distribuição adequada das tarefas. A aplicação deve ser intuitiva, responsiva e segura, garantindo que os dados dos usuários e das solicitações sejam protegidos e, acima de tudo, que todos os chamados sejam atendidos de forma rápida e eficiente.

O nome do sistema é: Sistema de Controle de Chamados Internos.

# Decision:

Para isso, a equipe de desenvolvimento da Codifica decidiu adotar uma arquitetura orientada a eventos, utilizando microsserviços e aplicando conceitos de Domain Driven Design (DDD). A principal motivação é garantir que cada ação sobre um chamado seja registrada como um evento, tornando todas as ações realizadas no sistema rastreáveis e auditáveis. Além disso, a arquitetura permite que diferentes partes do sistema evoluam de forma independente.

A equipe dividiu o sistema em dois bounded contexts principais:

- **Contexto de Chamados:** Responsável por gerenciar o ciclo de vida dos chamados, desde a criação até a resolução, incluindo a atribuição de tarefas, acompanhamento do status e prioridade. Este contexto utilizará **Event Sourcing**: o histórico de eventos de cada chamado (`ChamadoAberto`, `ChamadoAtribuido`, `PrioridadeAlterada`, `RespostaAdicionada`, `ChamadoResolvido`, entre outros) é a fonte da verdade, e o estado atual é reconstruído a partir dele. As consultas (como "meus chamados" e "fila do suporte") serão atendidas por projeções atualizadas a partir desses eventos.
- **Contexto de Funcionários:** Focado na gestão dos usuários do sistema, incluindo administradores, equipe de suporte e usuários finais. Este contexto será responsável por autenticação, autorização, gerenciamento de perfis e permissões. Utilizará persistência convencional e publicará eventos quando houver mudanças relevantes para outros contextos (`FuncionarioCadastrado`, `PapelAlterado`, `FuncionarioDesativado`).

E um subcontexto genérico para lidar com eventos e integrações entre os dois bounded contexts principais:

- **Subcontexto de Integração e Eventos:** Responsável pela comunicação assíncrona entre os contextos por meio de mensageria, garantindo que as ações realizadas em um contexto sejam refletidas no outro. Cada serviço grava seus eventos e uma tabela de saída (outbox) na mesma transação, e um processo em segundo plano publica esses eventos no broker. Este subcontexto também lidará com integrações externas, como envio de notificações por e-mail ou integração com sistemas de terceiros no futuro.

## Stack tecnológica:

- **Backend:** Go, um serviço por bounded context (`calls-service` e `employees-service`).
- **Frontend:** Vue com Vite.
- **Banco de dados:** PostgreSQL, com o event store do Contexto de Chamados implementado como uma tabela de eventos.
- **Mensageria:** Apache Kafka.
- **Infraestrutura:** Docker, Kubernetes com Helm, em ambiente Linux. Docker Compose para o ambiente local de desenvolvimento.
- **CI/CD:** pipeline automatizado para testes, build das imagens e deploy.

# Alternatives Considered:

- **Monolito Modular:** Uma única aplicação com os mesmos contextos separados em módulos, comunicando-se por eventos internos. Teria menor custo de infraestrutura e operação, mantendo a rastreabilidade. Foi descartado porque a equipe prioriza a implantação independente de cada contexto e a facilidade de incorporar integrações futuras por meio da mensageria. Caso esses requisitos percam relevância, esta alternativa deve ser reavaliada.
- **Arquitetura em Camadas:** Embora seja uma abordagem tradicional e bem compreendida, a organização apenas por camadas técnicas (apresentação, negócio e dados) tende a misturar as regras de chamados e de funcionários nas mesmas camadas, dificultando a evolução independente. Além disso, a persistência apenas do estado atual não atende à rastreabilidade sem uma trilha de auditoria paralela. A organização em camadas continuará presente dentro de cada serviço (domínio, aplicação e infraestrutura).
- **Arquitetura Orientada a Serviços (SOA):** Poderia oferecer modularidade e reutilização de componentes, mas depende de um barramento central que concentra a lógica de integração e se torna um ponto único de acoplamento e de falha, o que é desproporcional ao porte do sistema.

# Trade-offs:

## Benefícios:
- **Rastreabilidade e auditoria:** com Event Sourcing, o histórico completo de cada chamado pode ser consultado e o seu estado reconstruído a qualquer momento.
- **Tolerância a falhas:** a comunicação assíncrona permite que o sistema continue funcionando mesmo quando alguns componentes falham. Se o envio de notificações estiver indisponível, a abertura e o atendimento de chamados continuam, e os eventos são processados quando o componente volta.
- **Escalabilidade:** cada serviço pode ser escalado de forma independente no Kubernetes, de acordo com a sua demanda.
- **Flexibilidade:** novas integrações podem ser adicionadas como novos consumidores de eventos, sem alterar os serviços existentes.

## Custos:
- **Complexidade:** a arquitetura exige habilidades específicas da equipe e aumenta o tempo necessário para implementar novas funcionalidades.
- **Consistência eventual:** um contexto pode enxergar dados do outro com pequeno atraso. Por exemplo, um chamado pode ser atribuído a um atendente que acabou de ser desativado, e as regras de negócio precisam tratar esse caso.
- **Eventos duplicados:** os consumidores precisam processar o mesmo evento mais de uma vez sem gerar efeitos repetidos.
- **Custo geral:** Kubernetes, mensageria e múltiplos serviços aumentam os custos de desenvolvimento, operação e monitoramento.

# Consequences:

Como consequência da adoção dessa arquitetura, a equipe de desenvolvimento precisará investir em treinamento e capacitação em Domain Driven Design, Event Sourcing e arquitetura orientada a eventos. Será necessário implementar desde o início mecanismos de monitoramento e rastreamento entre serviços, para que o sistema funcione de forma eficiente e confiável mesmo em situações de falha ou alta demanda. Decisões derivadas desta, como a estratégia de versionamento dos eventos e a autenticação entre serviços, serão registradas em ADRs próprios.

# Security Considerations:

A arquitetura requer atenção especial à proteção dos dados dos usuários e das solicitações, à autenticação e autorização adequadas e aos mecanismos de auditoria. A autenticação será centralizada no Contexto de Funcionários, e cada serviço validará o token recebido e aplicará as suas próprias regras de autorização.

Como eventos armazenados não podem ser apagados, os eventos do Contexto de Chamados guardarão apenas o identificador dos funcionários, mantendo os dados pessoais no Contexto de Funcionários, onde podem ser corrigidos ou removidos, em conformidade com a LGPD.

Além disso, serão implementadas medidas contra ataques externos, como injeção de SQL, cross-site scripting (XSS) e negação de serviço (DoS), com uso obrigatório de HTTPS e acesso ao broker restrito aos serviços autorizados.