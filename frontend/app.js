const apiBase = window.location.protocol === "file:" ? "http://localhost:8080/api" : "/api";

const state = {
  usuarioID: Number(localStorage.getItem("usuarioID")) || null,
  categorias: [],
  orcamentos: [],
  categorySummaries: [],
  monthlySummary: null,
  monthlyBudgetTotal: 0,
};

const screens = document.querySelectorAll(".screen");
const navItems = document.querySelectorAll(".nav-item");
const appToast = document.getElementById("appToast");
const toast = window.bootstrap ? new window.bootstrap.Toast(appToast, { delay: 1400 }) : null;
const editCategoriaModalElement = document.getElementById("editCategoriaModal");
const editCategoriaModal = window.bootstrap ? new window.bootstrap.Modal(editCategoriaModalElement) : null;

function showToast(message) {
  appToast.querySelector(".toast-body").textContent = message;
  if (toast) {
    toast.show();
    return;
  }
  alert(message);
}

function setScreen(screenName) {
  screens.forEach((screen) => {
    screen.classList.toggle("active", screen.dataset.screen === screenName);
  });

  navItems.forEach((item) => {
    item.classList.toggle("active", item.dataset.target === screenName);
  });

  if (screenName === "categorias" || screenName === "transacao") {
    loadCategorias();
  }

  if (screenName === "relatorio") {
    loadMonthlySummary();
  }

  if (screenName === "projecao") {
    loadProjection();
  }
}

async function request(path, options = {}) {
  const response = await fetch(`${apiBase}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  const payload = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(payload.message || "Erro ao chamar API");
  }

  return payload.data;
}

function getUsuarioID() {
  if (state.usuarioID) {
    return state.usuarioID;
  }

  const typedID = Number(prompt("Informe o ID do usuário para continuar"));
  if (!typedID) {
    throw new Error("Usuário não informado");
  }

  state.usuarioID = typedID;
  localStorage.setItem("usuarioID", String(typedID));
  return typedID;
}

function getFormElement(id) {
  const formElement = document.getElementById(id);
  if (!(formElement instanceof HTMLFormElement)) {
    throw new Error(`Formulario "${id}" nao encontrado`);
  }

  return formElement;
}

async function loadCategorias() {
  let usuarioID;
  try {
    usuarioID = getUsuarioID();
  } catch {
    return;
  }

  try {
    state.categorias = await request(`/usuarios/${usuarioID}/categorias`);
    await Promise.all([loadOrcamentos(), loadCategorySummaries()]);
    renderCategorias();
    renderCategoriaOptions();
  } catch (error) {
    showToast(error.message);
  }
}

function renderCategoriaOptions() {
  const select = document.getElementById("categoria");
  select.innerHTML = "";

  state.categorias.forEach((categoria) => {
    const option = document.createElement("option");
    option.value = categoria.id;
    option.textContent = categoria.nome;
    select.appendChild(option);
  });
}

function renderCategorias() {
  const list = document.getElementById("categoriasList");
  list.innerHTML = "";

  if (!state.categorias.length) {
    list.innerHTML = '<div class="empty-state">Nenhuma categoria encontrada</div>';
    return;
  }

  state.categorias.forEach((categoria) => {
    const orcamento = findOrcamentoByCategoriaID(categoria.id);
    const summary = findCategorySummaryByCategoriaID(categoria.id);
    const gasto = Number(summary?.gasto || 0);
    const disponivel = Number(summary?.disponivel ?? (orcamento?.limite || 0));
    const percentualUtilizado = Number(summary?.percentual_utilizado || 0);
    const progressWidth = orcamento?.limite > 0 ? Math.min(percentualUtilizado, 100) : 0;
    const statusClass = disponivel < 0 ? "is-over" : "is-safe";
    const disponivelLabel = disponivel < 0 ? "Excedido" : "Disponível";
    const row = document.createElement("div");
    row.className = "category-row";
    row.innerHTML = `
      <div>
        <div class="category-name">${categoria.nome}</div>
        <div class="category-budget">Orçamento: ${formatMoney(orcamento?.limite || 0)}</div>
        <div class="category-metrics">
          <span>Gasto: ${formatMoney(gasto)}</span>
          <span class="${statusClass}">${disponivelLabel}: ${formatMoney(Math.abs(disponivel))}</span>
        </div>
        <div class="progress category-progress" role="progressbar" aria-label="Consumo do orçamento da categoria">
          <div class="progress-bar ${disponivel < 0 ? "bg-danger" : "bg-success"}" style="width: ${progressWidth}%"></div>
        </div>
      </div>
      <div class="category-actions">
        <button class="btn btn-outline-primary" type="button" data-budget="${categoria.id}">Limite</button>
        <button class="btn btn-primary" type="button" data-edit="${categoria.id}">Editar</button>
        <button class="btn btn-danger" type="button" data-delete="${categoria.id}">Excluir</button>
      </div>
    `;
    list.appendChild(row);
  });
}

async function loadOrcamentos(mes = getMonthInputValue("mesCategorias")) {
  if (!state.usuarioID || !mes) {
    state.orcamentos = [];
    return;
  }

  state.orcamentos = await request(`/usuarios/${state.usuarioID}/orcamentos?mes=${monthToDateParam(mes)}`);
}

function findOrcamentoByCategoriaID(categoriaID) {
  return state.orcamentos.find((orcamento) => orcamento.categoria_id === categoriaID);
}

async function loadCategorySummaries(mes = getMonthInputValue("mesCategorias")) {
  if (!state.usuarioID || !mes) {
    state.categorySummaries = [];
    renderCategoryOverview();
    return;
  }

  state.categorySummaries = await request(`/usuarios/${state.usuarioID}/relatorios/categorias?mes=${monthToDateParam(mes)}`);
  renderCategoryOverview();
}

function findCategorySummaryByCategoriaID(categoriaID) {
  return state.categorySummaries.find((summary) => summary.categoria_id === categoriaID);
}

function renderCategoryOverview() {
  const totals = state.categorySummaries.reduce((accumulator, item) => {
    accumulator.orcado += Number(item.orcamento || 0);
    accumulator.gasto += Number(item.gasto || 0);
    accumulator.disponivel += Number(item.disponivel || 0);
    return accumulator;
  }, { orcado: 0, gasto: 0, disponivel: 0 });

  document.getElementById("categoriasTotalOrcado").textContent = formatMoney(totals.orcado);
  document.getElementById("categoriasTotalGasto").textContent = formatMoney(totals.gasto);
  document.getElementById("categoriasTotalDisponivel").textContent = formatMoney(totals.disponivel);
}

document.addEventListener("click", async (event) => {
  const goButton = event.target.closest("[data-action='go']");
  if (goButton) {
    setScreen(goButton.dataset.target);
    return;
  }

  const editButton = event.target.closest("[data-edit]");
  if (editButton) {
    const categoria = state.categorias.find((item) => item.id === Number(editButton.dataset.edit));
    if (!categoria) {
      return;
    }

    if (!editCategoriaModal) {
      const nome = prompt("Novo nome da categoria", categoria.nome);
      if (!nome) {
        return;
      }

      await updateCategoria(categoria.id, nome);
      return;
    }

    document.getElementById("editCategoriaID").value = categoria.id;
    document.getElementById("editCategoriaNome").value = categoria.nome;
    editCategoriaModal.show();
    return;
  }

  const budgetButton = event.target.closest("[data-budget]");
  if (budgetButton) {
    const categoriaID = Number(budgetButton.dataset.budget);
    const categoria = state.categorias.find((item) => item.id === categoriaID);
    const orcamento = findOrcamentoByCategoriaID(categoriaID);
    const limite = prompt(`Quanto quer gastar em ${categoria?.nome || "esta categoria"}?`, orcamento?.limite || "");
    if (!limite) {
      return;
    }

    try {
      await saveOrcamento(categoriaID, parseAmount(limite), orcamento);
      showToast("Orçamento salvo");
      await loadCategorias();
    } catch (error) {
      showToast(error.message);
    }
    return;
  }

  const deleteButton = event.target.closest("[data-delete]");
  if (deleteButton) {
    const id = Number(deleteButton.dataset.delete);
    if (!confirm("Excluir esta categoria?")) {
      return;
    }

    try {
      await request(`/categorias/${id}`, { method: "DELETE" });
      showToast("Categoria excluída");
      await loadCategorias();
    } catch (error) {
      showToast(error.message);
    }
  }
});

document.getElementById("mesCategorias").addEventListener("change", loadCategorias);
document.getElementById("mesRelatorio").addEventListener("change", loadMonthlySummary);

const usuarioForm = getFormElement("usuarioForm");
const transacaoForm = getFormElement("transacaoForm");
const categoriaForm = getFormElement("categoriaForm");

usuarioForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const formElement = usuarioForm;
  const form = new FormData(formElement);

  try {
    const usuario = await request("/usuarios", {
      method: "POST",
      body: JSON.stringify({
        nome: form.get("nome"),
        email: form.get("email"),
      }),
    });

    state.usuarioID = usuario.id;
    localStorage.setItem("usuarioID", String(usuario.id));
    showToast("Usuário criado");
    formElement.reset();
    setScreen("transacao");
  } catch (error) {
    showToast(error.message);
  }
});

transacaoForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const formElement = transacaoForm;
  const form = new FormData(formElement);

  try {
    await request("/transacoes", {
      method: "POST",
      body: JSON.stringify({
        usuario_id: getUsuarioID(),
        categoria_id: Number(form.get("categoria_id")),
        valor: Number(form.get("valor")),
        data: `${form.get("data")}T00:00:00Z`,
        descricao: form.get("descricao"),
        tipo: form.get("tipo"),
        parcelas: Number(form.get("parcelas")) || 1,
      }),
    });

    showToast("Transação criada");
    formElement.reset();
    document.getElementById("parcelas").value = "1";
    setToday();
  } catch (error) {
    showToast(error.message);
  }
});

categoriaForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const formElement = categoriaForm;
  const form = new FormData(formElement);

  try {
    await request("/categorias", {
      method: "POST",
      body: JSON.stringify({
        usuario_id: getUsuarioID(),
        nome: form.get("nome"),
      }),
    });

    showToast("Categoria criada");
    formElement.reset();
    await loadCategorias();
  } catch (error) {
    showToast(error.message);
  }
});

document.getElementById("editCategoriaForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  const id = Number(document.getElementById("editCategoriaID").value);
  const nome = document.getElementById("editCategoriaNome").value;

  try {
    await updateCategoria(id, nome);
    if (editCategoriaModal) {
      editCategoriaModal.hide();
    }
  } catch (error) {
    showToast(error.message);
  }
});

async function updateCategoria(id, nome) {
  await request(`/categorias/${id}`, {
    method: "PUT",
    body: JSON.stringify({
      id,
      nome,
      usuario_id: getUsuarioID(),
    }),
  });

  showToast("Categoria atualizada");
  await loadCategorias();
}

async function saveOrcamento(categoriaID, limite, orcamento) {
  if (!limite || limite <= 0) {
    throw new Error("Informe um valor maior que zero");
  }

  const body = {
    usuario_id: getUsuarioID(),
    categoria_id: categoriaID,
    limite,
    mes: `${getMonthInputValue("mesCategorias")}-01T00:00:00Z`,
  };

  if (orcamento) {
    await request(`/orcamentos/${orcamento.id}`, {
      method: "PUT",
      body: JSON.stringify({ id: orcamento.id, ...body }),
    });
    return;
  }

  await request("/orcamentos", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

async function loadMonthlySummary() {
  let usuarioID;
  try {
    usuarioID = getUsuarioID();
  } catch {
    return;
  }

  try {
    const mes = getMonthInputValue("mesRelatorio");
    const summary = await request(`/usuarios/${usuarioID}/relatorios/mensal?mes=${monthToDateParam(mes)}`);
    state.monthlySummary = summary;
    renderMonthlySummary(summary);
  } catch (error) {
    showToast(error.message);
  }
}

function renderMonthlySummary(summary) {
  const receita = Number(summary.total_receita || 0);
  const despesa = Number(summary.total_despesa || 0);
  const saldo = Number(summary.saldo_atual || 0);
  const max = Math.max(receita, despesa, Math.abs(saldo), 1);

  setBar("receitaBar", "receitaBarValue", receita, max);
  setBar("despesaBar", "despesaBarValue", despesa, max);
  setBar("saldoBar", "saldoBarValue", Math.max(saldo, 0), max);

  document.getElementById("totalReceitas").textContent = formatMoney(receita);
  document.getElementById("totalDespesas").textContent = formatMoney(despesa);
  document.getElementById("saldoAtual").textContent = formatMoney(saldo);
}

function setBar(barID, valueID, value, max) {
  document.getElementById(barID).style.height = `${Math.max((value / max) * 145, value > 0 ? 8 : 0)}px`;
  document.getElementById(valueID).textContent = formatMoney(value);
}

async function loadProjection() {
  let usuarioID;
  try {
    usuarioID = getUsuarioID();
  } catch {
    return;
  }

  try {
    const mes = currentMonthValue();
    const projection = await request(`/usuarios/${usuarioID}/projecao/comprometimento?mes=${monthToDateParam(mes)}&meses=4`);
    renderProjection(projection);
  } catch (error) {
    showToast(error.message);
  }
}

function renderProjection(projection) {
  const list = document.getElementById("commitmentList");
  list.innerHTML = "";

  if (!projection.length) {
    list.innerHTML = '<div class="empty-state">Nenhuma parcela futura encontrada</div>';
    return;
  }

  projection.forEach((item) => {
    const receita = Number(item.total_receita || 0);
    const despesa = Number(item.total_despesa || 0);
    const saldo = Number(item.saldo_projetado || 0);
    const percentual = Number(item.percentual_comprometido || 0);
    const progressWidth = receita > 0 ? Math.min(percentual, 100) : 0;

    list.insertAdjacentHTML("beforeend", `
      <div class="commitment-row">
        <div class="commitment-header">
          <strong>${formatMonth(item.mes)}</strong>
          <span>${Math.round(percentual)}%</span>
        </div>
        <div class="progress commitment-progress" role="progressbar" aria-label="Renda comprometida">
          <div class="progress-bar bg-danger" style="width: ${progressWidth}%"></div>
        </div>
        <div class="commitment-values">
          <span>Parcelas: ${formatMoney(despesa)}</span>
          <span>Receita: ${formatMoney(receita)}</span>
          <span>Saldo: ${formatMoney(saldo)}</span>
        </div>
      </div>
    `);
  });
}

function formatMoney(value) {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
    maximumFractionDigits: 0,
  }).format(Number(value || 0));
}

function formatMonth(value) {
  const date = new Date(value);
  return new Intl.DateTimeFormat("pt-BR", {
    month: "short",
    year: "numeric",
    timeZone: "UTC",
  }).format(date);
}

function parseAmount(value) {
  const normalizedValue = String(value).trim();
  if (normalizedValue.includes(",")) {
    return Number(normalizedValue.replaceAll(".", "").replace(",", "."));
  }

  return Number(normalizedValue);
}

function getMonthInputValue(inputID) {
  const input = document.getElementById(inputID);
  if (!input.value) {
    input.value = currentMonthValue();
  }
  return input.value;
}

function monthToDateParam(monthValue) {
  return `${monthValue}-01`;
}

function currentMonthValue() {
  return new Date().toISOString().slice(0, 7);
}

function setToday() {
  document.getElementById("data").value = new Date().toISOString().slice(0, 10);
}

["mesCategorias", "mesRelatorio"].forEach((inputID) => {
  document.getElementById(inputID).value = currentMonthValue();
});
setToday();
setScreen("usuario");
