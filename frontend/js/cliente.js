async function init() {
    try {
        const data = await fetchClientes();
        renderTable(data);
    } catch (error) {
        console.error("Erro ao carregar clientes:", error);
        document.querySelector('#clientesTable tbody').innerHTML =
            '<tr><td colspan="3">Erro ao carregar dados. Verifique se o servidor está rodando.</td></tr>';
    }
}

// Variável global para armazenar os clientes
let clientes = [];

// Função para carregar clientes do servidor
async function loadClientes() {
    try {
        clientes = await fetchClientes();
        renderTable(clientes);
    } catch (error) {
        console.error("Erro ao carregar clientes:", error);
    }
}

function renderTable(data) {
    const tbody = document.querySelector('#clientesTable tbody');
    tbody.innerHTML = '';

    if (data.length === 0) {
        tbody.innerHTML = '<tr><td colspan="3">Nenhum cliente encontrado</td></tr>';
        return;
    }

    data.forEach(c => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
      <td>${c.nome}</td>
      <td>${c.telefone}</td>
      <td>
        <button class="action-btn edit-btn" data-id="${c.id}">Editar</button>
        <button class="action-btn delete-btn" data-id="${c.id}">Deletar</button>
      </td>`;
        tbody.appendChild(tr);
    });
    attachEventListeners();
}

// Liga eventos de editar e deletar
function attachEventListeners() {
    document.querySelectorAll('.edit-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            window.location.href = `/partials/cliente_form.html?id=${btn.dataset.id}`;
        });
    });
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', async () => {
            if (confirm('Confirma a exclusão?')) {
                try {
                    await deleteCliente(btn.dataset.id);
                    loadClientes(); // Corrigido: descomentado para recarregar a lista após deletar
                } catch (error) {
                    console.error("Erro ao deletar cliente:", error);
                    alert("Erro ao deletar cliente. Verifique se o servidor está rodando.");
                }
            }
        });
    });
}

// Inicialização ao carregar a página
document.addEventListener('DOMContentLoaded', () => {
    // Carrega o sidebar usando caminho relativo
    fetch('/partials/sidebar.html')
        .then(res => {
            if (!res.ok) {
                throw new Error(`Erro ao carregar sidebar: ${res.status}`);
            }
            return res.text();
        })
        .then(html => {
            document.getElementById('sidebar').innerHTML = html;
        })
        .catch(error => {
            console.error("Erro ao carregar sidebar:", error);
            document.getElementById('sidebar').innerHTML = '<p>Menu não disponível</p>';
        });

    // Configura busca dinâmica
    document.getElementById('searchInput').addEventListener('input', e => {
        const term = e.target.value.toLowerCase();
        const filtered = clientes.filter(c =>
            c.nome.toLowerCase().includes(term) ||
            c.telefone.includes(term)
        );
        renderTable(filtered);
    });

    // Botão Novo Cliente
    document.getElementById('newBtn').addEventListener('click', () => {
        window.location.href = '/partials/cliente_form.html';
    });

    loadClientes(); // Use a função de carregar clientes
});