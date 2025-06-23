const API_BASE_URL = 'http://localhost:8080/api/v1';

// 响应处理函数
const handleResponse = async (response) => {
    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `请求失败，状态码: ${response.status}`);
    }
    return response.json();
};

export const api = {
    // 获取题目列表
    getQuestions: async ({ type, keyword, pageNum, pageSize }) => {
        const params = new URLSearchParams();
        if (type !== 0) params.append('type', type);
        if (keyword) params.append('keyword', keyword);
        params.append('pageNum', pageNum);
        params.append('pageSize', pageSize);
        console.log(type);
        console.log(keyword);
        console.log(pageNum);
        console.log(pageSize);

        const response = await fetch(`${API_BASE_URL}/questions?${params}`);
        return handleResponse(response);
    },

    // 获取单个题目详情
    getQuestion: async (id) => {
        const response = await fetch(`${API_BASE_URL}/questions/${id}`);
        return handleResponse(response);
    },

    // 创建题目
    createQuestion: async (question) => {
        const response = await fetch(`${API_BASE_URL}/questions`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(question)
        });
        return handleResponse(response);
    },

    // 批量创建题目
    batchCreateQuestions: async (questions) => {
        const response = await fetch(`${API_BASE_URL}/questions/batch`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(questions)
        });
        return handleResponse(response);
    },

    // 更新题目
    updateQuestion: async (id, question) => {
        const response = await fetch(`${API_BASE_URL}/questions/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(question)
        });
        return handleResponse(response);
    },

    // 删除题目
    deleteQuestions: async (ids) => {
        // 确保ids始终是数组
        const idArray = Array.isArray(ids) ? ids : [ids];

        const response = await fetch(`${API_BASE_URL}/questions`, {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ids: idArray.map(String) }) // 确保ID是字符串
        });
        return handleResponse(response);
    },

    // AI生成题目
    generateQuestions: async ({ type, language, amount }) => {
        const response = await fetch(`${API_BASE_URL}/questions/ai-generate`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                type,
                language,
                amount
            })
        });
        return handleResponse(response);
    }
};