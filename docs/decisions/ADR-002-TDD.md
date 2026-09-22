# ADR-002: TDD (Test-Driven Development)

status: Proposto
data: 2026-09-21

# Context:

Para o desenvolvimento do Sistema de Controle de Chamados Internos, foi adotada a abordagem de TDD (Test-Driven Development) como prática de desenvolvimento. O TDD é uma metodologia que prioriza a escrita de testes antes da implementação do código, garantindo que o software seja desenvolvido com foco na qualidade e na funcionalidade desejada.

# Decision:

A prática de TDD será aplicada em todas as fases do desenvolvimento, desde a criação de novos recursos até a refatoração de código existente. A equipe seguirá o ciclo de TDD: escrever um teste que falhe, implementar o código necessário para passar no teste e, em seguida, refatorar o código mantendo os testes verdes.

# Alternatives Considered:

- **Desenvolvimento sem testes:** A equipe poderia optar por desenvolver o sistema sem a prática de TDD, confiando apenas em testes manuais ou testes automatizados escritos após a implementação. No entanto, essa abordagem aumenta o risco de introdução de bugs e reduz a qualidade do código.
- **Testes automatizados após a implementação:** Outra alternativa seria escrever testes automatizados apenas após a implementação do código. Embora isso possa fornecer alguma cobertura de testes, não oferece o mesmo nível de garantia de qualidade e feedback rápido que o TDD proporciona.

# Trade-offs:

## Benefícios:
- **Maior qualidade do código:** A prática de TDD promove um design mais robusto e modular, facilitando a manutenção e evolução do sistema.
- **Redução de bugs:** A execução frequente de testes ajuda a identificar e corrigir problemas rapidamente, reduzindo a quantidade de bugs em produção.
- **Documentação viva:** Os testes servem como documentação do comportamento esperado do sistema, facilitando a compreensão do código por novos membros da equipe.
- **Feedback rápido:** A execução frequente dos testes permite que a equipe receba feedback imediato sobre o impacto das mudanças no código, facilitando a identificação de regressões e problemas de integração.

## Custos:
- **Curva de aprendizado:** A equipe precisará investir tempo para se familiarizar com a abordagem de TDD e as ferramentas de teste utilizadas.
- **Tempo de desenvolvimento inicial:** A escrita de testes antes da implementação pode aumentar o tempo de desenvolvimento inicial, embora isso seja compensado pela redução de tempo gasto na correção de bugs e manutenção do código a longo prazo.
- **Manutenção dos testes:** À medida que o sistema evolui, os testes precisarão ser atualizados para refletir as mudanças no código, o que pode exigir esforço adicional da equipe.

# Consequences:

Como consequencia da adoção do TDD, se faz necessário o comprometimento em seguir rigorosamente a prática de escrever testes antes da implementação do código. Inicialmente sua curva de aprendizado pode impactar a produtividade, mas a longo prazo trará benefícios significativos em termos de qualidade do código, redução de bugs e facilidade de manutenção.

# Security Considerations:

A prática de TDD contribui para a segurança do sistema, pois permite identificar vulnerabilidades e falhas de segurança durante o desenvolvimento, garantindo que o código seja testado e validado antes de ser implementado em produção.
