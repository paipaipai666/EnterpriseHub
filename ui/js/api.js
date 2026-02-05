const API_BASE_URL = 'http://localhost:8080/api/v1';

const APIConfig = {
    baseURL: API_BASE_URL,
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json',
    }
};

const UserAPI = {
    async register(username, password) {
        try {
            const response = await fetch(`${API_BASE_URL}/users/register`, {
                method: 'POST',
                headers: APIConfig.headers,
                body: JSON.stringify({ username, password })
            });
            return await response.json();
        } catch (error) {
            console.error('User register error:', error);
            throw error;
        }
    },

    async getById(id) {
        try {
            const response = await fetch(`${API_BASE_URL}/users/get/${id}`, {
                method: 'GET',
                headers: this.getAuthHeaders()
            });
            return await response.json();
        } catch (error) {
            console.error('Get user error:', error);
            throw error;
        }
    },

    async getAll() {
        try {
            const response = await fetch(`${API_BASE_URL}/users/get_all`, {
                method: 'GET',
                headers: this.getAuthHeaders()
            });
            return await response.json();
        } catch (error) {
            console.error('Get all users error:', error);
            throw error;
        }
    },

    async getByUsername(username) {
        try {
            const response = await fetch(`${API_BASE_URL}/users/get_by_username/${username}`, {
                method: 'GET',
                headers: this.getAuthHeaders()
            });
            return await response.json();
        } catch (error) {
            console.error('Get user by username error:', error);
            throw error;
        }
    },

    getAuthHeaders() {
        const token = localStorage.getItem('token');
        return {
            'Content-Type': 'application/json',
            ...(token && { 'Authorization': `Bearer ${token}` })
        };
    }
};

const AuthAPI = {
    async login(username, password) {
        try {
            const response = await fetch(`${API_BASE_URL}/auth/login`, {
                method: 'POST',
                headers: APIConfig.headers,
                body: JSON.stringify({ username, password })
            });
            return await response.json();
        } catch (error) {
            console.error('Login error:', error);
            throw error;
        }
    },

    async logout() {
        try {
            const response = await fetch(`${API_BASE_URL}/auth/logout`, {
                method: 'DELETE',
                headers: UserAPI.getAuthHeaders()
            });
            return await response.json();
        } catch (error) {
            console.error('Logout error:', error);
            throw error;
        }
    }
};

const OrderAPI = {
    async create(userId, amount) {
        try {
            const response = await fetch(`${API_BASE_URL}/order/create`, {
                method: 'POST',
                headers: UserAPI.getAuthHeaders(),
                body: JSON.stringify({ userId, amount: parseFloat(amount) })
            });
            return await response.json();
        } catch (error) {
            console.error('Create order error:', error);
            throw error;
        }
    },

    async getById(id) {
        try {
            const response = await fetch(`${API_BASE_URL}/order/get/${id}`, {
                method: 'GET',
                headers: UserAPI.getAuthHeaders()
            });
            return await response.json();
        } catch (error) {
            console.error('Get order error:', error);
            throw error;
        }
    },

    async getList(userId) {
        try {
            const response = await fetch(`${API_BASE_URL}/order/list/${userId}`, {
                method: 'GET',
                headers: UserAPI.getAuthHeaders()
            });
            return await response.json();
        } catch (error) {
            console.error('Get order list error:', error);
            throw error;
        }
    },

    async cancel(id) {
        try {
            const response = await fetch(`${API_BASE_URL}/order/cancel/${id}`, {
                method: 'PUT',
                headers: UserAPI.getAuthHeaders()
            });
            return await response.json();
        } catch (error) {
            console.error('Cancel order error:', error);
            throw error;
        }
    },

    async pay(id, method) {
        try {
            const response = await fetch(`${API_BASE_URL}/order/pay/${id}`, {
                method: 'POST',
                headers: UserAPI.getAuthHeaders(),
                body: JSON.stringify({ method: parseInt(method) })
            });
            return await response.json();
        } catch (error) {
            console.error('Pay order error:', error);
            throw error;
        }
    }
};

const PaymentMethod = {
    UNKNOWN: 0,
    ALIPAY: 1,
    WECHAT: 2,
    BANK_CARD: 3
};

const PaymentMethodText = {
    0: '未知',
    1: '支付宝',
    2: '微信支付',
    3: '银行卡'
};

const OrderStatus = {
    CREATED: 'CREATED',
    PENDING: 'PENDING',
    PAID: 'PAID',
    CANCELLED: 'CANCELLED'
};

const OrderStatusText = {
    CREATED: '待支付',
    PENDING: '支付中',
    PAID: '已支付',
    CANCELLED: '已取消',
    created: '待支付',
    pending: '支付中',
    paid: '已支付',
    cancelled: '已取消'
};

function formatDate(dateValue) {
    if (!dateValue) return '-';
    let date;
    if (typeof dateValue === 'number') {
        date = new Date(dateValue * 1000);
    } else if (typeof dateValue === 'string') {
        date = new Date(dateValue);
    } else {
        return '-';
    }
    if (isNaN(date.getTime())) return '-';
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function formatCurrency(amount) {
    return '¥' + parseFloat(amount).toFixed(2);
}

function showToast(message, type = 'info') {
    const existingToast = document.querySelector('.toast-container');
    if (existingToast) existingToast.remove();

    const container = document.createElement('div');
    container.className = 'fixed top-4 right-4 z-50 toast-container space-y-2';

    const toast = document.createElement('div');
    const bgColor = type === 'success' ? 'bg-green-500' : type === 'error' ? 'bg-red-500' : 'bg-blue-500';
    const icon = type === 'success' ? 'check-circle' : type === 'error' ? 'exclamation-circle' : 'info-circle';

    toast.className = `${bgColor} text-white px-6 py-3 rounded-lg shadow-lg flex items-center transform transition-all duration-300 translate-x-full opacity-0`;
    toast.innerHTML = `
        <i class="fas fa-${icon} mr-2"></i>
        <span>${message}</span>
    `;

    container.appendChild(toast);
    document.body.appendChild(container);

    setTimeout(() => {
        toast.classList.remove('translate-x-full', 'opacity-0');
    }, 10);

    setTimeout(() => {
        toast.classList.add('translate-x-full', 'opacity-0');
        setTimeout(() => {
            container.remove();
        }, 300);
    }, 3000);
}

function showLoading(containerId) {
    const container = document.getElementById(containerId);
    if (container) {
        container.innerHTML = `
            <div class="flex items-center justify-center py-12">
                <i class="fas fa-spinner fa-spin text-3xl text-primary-500"></i>
                <span class="ml-3 text-gray-500">加载中...</span>
            </div>
        `;
    }
}

function showEmpty(containerId, message = '暂无数据') {
    const container = document.getElementById(containerId);
    if (container) {
        container.innerHTML = `
            <div class="flex flex-col items-center justify-center py-12 text-gray-400">
                <i class="fas fa-inbox text-5xl mb-3"></i>
                <span>${message}</span>
            </div>
        `;
    }
}

function showError(containerId, message = '加载失败，请稍后重试') {
    const container = document.getElementById(containerId);
    if (container) {
        container.innerHTML = `
            <div class="flex flex-col items-center justify-center py-12 text-red-500">
                <i class="fas fa-exclamation-triangle text-5xl mb-3"></i>
                <span>${message}</span>
            </div>
        `;
    }
}

function confirmDialog(message, onConfirm) {
    if (confirm(message)) {
        onConfirm();
    }
}

const Modal = {
    open(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.classList.remove('hidden');
            document.body.classList.add('overflow-hidden');
        }
    },

    close(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.classList.add('hidden');
            document.body.classList.remove('overflow-hidden');
        }
    },

    closeAll() {
        document.querySelectorAll('[id^="modal-"]').forEach(modal => {
            modal.classList.add('hidden');
        });
        document.body.classList.remove('overflow-hidden');
    }
};

document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') {
        Modal.closeAll();
    }
});

window.addEventListener('click', function(e) {
    if (e.target.classList.contains('modal-overlay')) {
        Modal.closeAll();
    }
});
