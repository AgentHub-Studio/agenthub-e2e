# agenthub-e2e — Testes End-to-End

Suíte de testes e2e do AgentHub cobrindo os fluxos críticos do sistema:
backend Java, serviços IA (embedding + extractor) e isolamento multi-tenant.

## Pré-requisitos

| Requisito | Versão mínima |
|---|---|
| Java | 21 |
| Maven | 3.9+ |
| Stack Docker | `docker compose --profile backend up -d` em `agenthub-infra/` |
| Keycloak | Rodando em `localhost:8080` (padrão do compose) |
| Backend | Rodando em `localhost:8081` |
| Embedding (opcional) | Rodando em `localhost:8092` (Fluxo 6) |
| Extractor (opcional) | Rodando em `localhost:8093` (Fluxo 6) |

Para subir o stack completo (incluindo embedding e extractor):

```bash
cd agenthub-infra
docker compose --profile backend --profile full up -d
```

## Variáveis de Ambiente

| Variável | Padrão | Descrição |
|---|---|---|
| `BACKEND_URL` | `http://localhost:8081` | URL base do backend |
| `KEYCLOAK_URL` | `http://localhost:8080` | URL externa do Keycloak |
| `KEYCLOAK_ADMIN_USER` | `admin` | Usuário admin do realm master |
| `KEYCLOAK_ADMIN_PASSWORD` | `@admin#` | Senha do admin Keycloak |
| `EMBEDDING_URL` | `http://localhost:8092` | URL do agenthub-embedding |
| `EXTRACTOR_URL` | `http://localhost:8093` | URL do agenthub-extractor |
| `E2E_TENANT_PREFIX` | `e2e-test` | Prefixo dos slugs de tenant criados nos testes |
| `E2E_USER_PASSWORD` | `E2eTestPass#1` | Senha dos usuários admin criados nos testes |

## Como Rodar

### Todos os fluxos

```bash
cd agenthub-e2e
mvn test -De2e.skip=false
```

### Um fluxo específico

```bash
# Fluxo 1 — Provisionamento de Tenant
mvn test -De2e.skip=false -Dtest=TenantProvisioningE2ETest

# Fluxo 2 — Skill + Tool CRUD
mvn test -De2e.skip=false -Dtest=SkillToolCrudE2ETest

# Fluxo 3 — RAG Pipeline
mvn test -De2e.skip=false -Dtest=RagPipelineE2ETest

# Fluxo 4 — Chat com Agent
mvn test -De2e.skip=false -Dtest=AgentChatE2ETest

# Fluxo 5 — Isolamento Multi-Tenant
mvn test -De2e.skip=false -Dtest=MultiTenantIsolationE2ETest

# Fluxo 6 — Smoke Tests IA
mvn test -De2e.skip=false -Dtest=AiServicesSmokeE2ETest
```

### Com variáveis personalizadas

```bash
BACKEND_URL=http://meu-backend:8081 \
KEYCLOAK_ADMIN_PASSWORD=minha-senha \
mvn test -De2e.skip=false
```

## Fluxos Cobertos

| # | Nome | Descrição |
|---|---|---|
| 1 | Provisionamento de Tenant | Criação de tenant, provisioning Keycloak, autenticação JWT |
| 2 | Skill + Tool CRUD | Ciclo completo create → read → list → delete com paginação |
| 3 | RAG Pipeline | Upload de PDF → extração → chunking → embedding → INDEXED |
| 4 | Chat com Agent | Criação de agent, streaming SSE, histórico persistido |
| 5 | Isolamento Multi-Tenant | Tenant A não vê recursos de Tenant B e vice-versa |
| 6 | Smoke Tests IA | Embedding retorna vetor 1024-dim; Extractor extrai texto de PDF |

## Interpretando o Relatório

Ao final de `mvn test`, o Surefire gera relatórios em `target/e2e-reports/`.

```
Tests run: 32, Failures: 1, Errors: 0, Skipped: 0
```

| Status | Significado |
|---|---|
| ✅ PASS | Fluxo funcionando corretamente |
| ❌ FAIL | Assertion falhou — verificar mensagem de erro para causa raiz |
| ⚠️ ERROR | Exceção inesperada — pode indicar serviço offline ou timeout |
| ⏭️ SKIP | Teste pulado (e2e.skip=true ou dependência anterior falhou) |

Para ver detalhes de uma falha específica:

```bash
cat target/e2e-reports/TEST-dev.cezar.agenthub.e2e.flows.RagPipelineE2ETest.xml
```

## Limpeza Manual de Recursos

Todos os recursos criados pelos testes usam o prefixo `e2e-` no nome.
Se um teste falhar antes do `@AfterAll`, você pode limpá-los manualmente:

```bash
# Listar realms de teste no Keycloak
curl -s -H "Authorization: Bearer $(curl -s -d 'grant_type=password&client_id=admin-cli&username=admin&password=@admin#' http://localhost:8080/realms/master/protocol/openid-connect/token | jq -r .access_token)" \
  http://localhost:8080/admin/realms | jq '.[].realm' | grep e2e

# Deletar realm manualmente
curl -X DELETE -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:8080/admin/realms/e2e-test-XXXXXX
```

## Notas de Implementação

- **Sem Testcontainers** — testes rodam contra o stack Docker real
- **Isolamento por fluxo** — cada `*E2ETest` cria e destroi seu próprio tenant em `@BeforeAll`/`@AfterAll`
- **Recursos com prefixo `e2e-`** — identificáveis e removíveis manualmente se necessário
- **Fluxo 3 (RAG)** requer embedding e extractor rodando (profile `full`) — pula automaticamente se timeout
- **Fluxo 4 (Chat SSE)** requer um modelo LLM configurado no agent (Ollama local por padrão)
