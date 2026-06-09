# Contribuindo para o Sistema Runner

Obrigado por considerar contribuir para o Sistema Runner! Este documento fornece diretrizes e instruções para contribuições.

## Como Contribuir

### Reportar Bugs

Bugs são rastreados via [GitHub Issues](../../issues). Ao reportar um bug, inclua:

- Título descritivo
- Descrição clara do comportamento esperado vs. observado
- Passos para reproduzir
- Exemplos específicos para demonstrar os passos
- Screenshots ou logs, se aplicável
- Seu ambiente (OS, Go version, etc.)

### Sugerir Melhorias

Sugestões de funcionalidades são bem-vindas! Abra uma issue com:

- Título descritivo da funcionalidade
- Descrição do caso de uso
- Exemplos de como funcionaria
- Por que essa melhoria seria útil

### Pull Requests

1. **Fork** o repositório
2. **Clone** seu fork localmente
3. **Crie uma branch** para sua feature:
   ```bash
   git checkout -b feat/sua-feature
   ```
4. **Faça commits** atômicos com mensagens claras:
   ```bash
   git commit -m "feat: descrição clara da mudança"
   ```
5. **Siga as convenções de commit** (Conventional Commits):
   - `feat:` para novas funcionalidades
   - `fix:` para correções de bugs
   - `docs:` para documentação
   - `test:` para testes
   - `refactor:` para refatoração de código
   - `chore:` para tarefas de manutenção

6. **Push** para sua branch:
   ```bash
   git push origin feat/sua-feature
   ```
7. **Abra um Pull Request** com descrição clara

## Padrões de Código

### Go

- Siga [Go Code Review Comments](https://go.dev/doc/effective_go)
- Implemente `go fmt` e `go vet` sem erros
- Adicione testes para novas funcionalidades
- Mantenha funções curtas e focadas

### Java

- Siga [Java Code Conventions](https://www.oracle.com/java/technologies/javase/codeconventions-indentation.html)
- Use formatação consistente
- Adicione testes unitários
- Documente com Javadoc

## Testes

- **Testes unitários** para lógica de negócio
- **Testes de integração** para APIs e fluxos completos
- Execute localmente antes de fazer PR:
  ```bash
  go test ./...
  ```

## Documentação

- Mantenha o README atualizado
- Adicione comentários para código complexo
- Documente decisões técnicas em ADRs (Architecture Decision Records)

## Processo de Revisão

1. Um ou mais mantenedores revisarão seu PR
2. Solicitações de mudanças podem ser feitas
3. Após aprovação, seu PR será mergeado
4. Será criada uma release/tag quando apropriado

## Dúvidas?

Abra uma [discussion](../../discussions) ou uma [issue](../../issues) com a tag `question`.

## Código de Conduta

Todos os contribuidores devem ser respeitosos e inclusivos. Discriminação, assédio ou comportamento tóxico não será tolerado.

---

Obrigado por contribuir! 🙏
