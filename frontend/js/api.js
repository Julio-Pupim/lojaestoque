/*const API_BASE = 'http://localhost:8080';

// Busca todos os clientes
async function fetchClientes() {
    const res = await fetch(`${API_BASE}/clientes`);
    return await res.json();
}

// Deleta um cliente pelo ID
async function deleteCliente(id) {
    await fetch(`${API_BASE}/clientes/${id}`, { method: 'DELETE' });
}
*/

// Configuração da URL base da API
const API_BASE = 'http://localhost:8080'; // Certifique-se que esta porta corresponde à do seu servidor

// Função para tratar erros nas requisições fetch
async function fetchWithErrorHandling(url, options = {}) {
    try {
        const response = await fetch(url, options);

        if (!response.ok) {
            throw new Error(`Erro na requisição: ${response.status} - ${response.statusText}`);
        }

        // Se a resposta estiver vazia ou não for JSON, retorne null
        const contentType = response.headers.get('content-type');
        if (!contentType || !contentType.includes('application/json')) {
            return response.status === 204 ? null : await response.text();
        }

        return await response.json();
    } catch (error) {
        console.error(`Erro na requisição para ${url}:`, error);
        throw error;
    }
}

// Busca todos os clientes
async function fetchClientes() {
    return await fetchWithErrorHandling(`${API_BASE}/clientes`);
}

// Busca um cliente pelo ID
async function fetchClienteById(id) {
    return await fetchWithErrorHandling(`${API_BASE}/clientes/${id}`);
}

// Cria um novo cliente
async function createCliente(cliente) {
    return await fetchWithErrorHandling(`${API_BASE}/clientes`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(cliente)
    });
}

// Atualiza um cliente existente
async function updateCliente(id, cliente) {
    return await fetchWithErrorHandling(`${API_BASE}/clientes/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(cliente)
    });
}

// Deleta um cliente pelo ID
async function deleteCliente(id) {
    return await fetchWithErrorHandling(`${API_BASE}/clientes/${id}`, {
        method: 'DELETE'
    });
}