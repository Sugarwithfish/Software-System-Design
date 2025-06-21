import React, {useEffect, useState} from 'react';
import './App.css';
import { api } from './api'

function App() {
  const [isNavCollapsed, setIsNavCollapsed] = useState(false);
  const [activeQuestionType, setActiveQuestionType] = useState('全部');
  const [searchQuery, setSearchQuery] = useState('');
  const [showAddQuestionDropdown, setShowAddQuestionDropdown] = useState(false);
  const [selectedQuestions, setSelectedQuestions] = useState([]);
  const [questions, setQuestions] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const [totalQuestions, setTotalQuestions] = useState(0);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [questionToDelete, setQuestionToDelete] = useState(null);
  const [showEditModal, setShowEditModal] = useState(false);
  const [editingQuestion, setEditingQuestion] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [aiGeneratedQuestions, setAiGeneratedQuestions] = useState([]);
  const [showAIPreview, setShowAIPreview] = useState(false);
  const QUESTION_TYPE_MAP = {
    1: '单选题',
    2: '多选题',
    3: '编程题'
  };

  const REVERSE_QUESTION_TYPE_MAP = {
    '单选题': 1,
    '多选题': 2,
    '编程题': 3
  };

  useEffect(() => {
    const fetchQuestions = async () => {
      setIsLoading(true);
      try {
        const data = await api.getQuestions({
          pageNum: currentPage,
          pageSize: itemsPerPage,
          type: activeQuestionType === '全部' ? '0' : REVERSE_QUESTION_TYPE_MAP[activeQuestionType],
          keyword: searchQuery
        });
        if(data.code === 0 && data.data.records === null){
          return;
        }
        console.log(data);
        const formattedQuestions = data.data.records.map(q => ({
          ...q,
          type: QUESTION_TYPE_MAP[q.type],
          title: q.title,
          content: q.content,
          options: JSON.parse(q.options),
          id: q.id,
          answer: JSON.parse(q.answer),
          selected: false,
        }));
        setQuestions(formattedQuestions);
        setTotalQuestions(data.total);
      } catch (error) {
        console.error('加载题目失败:', error);
        alert('加载题目失败，请稍后重试');
      } finally {
        setIsLoading(false);
      }
    };

    fetchQuestions();
  }, [currentPage, itemsPerPage, activeQuestionType]);

  const toggleNav = () => {
    setIsNavCollapsed(!isNavCollapsed);
  };

  const handleQuestionTypeClick = (type) => {
    setActiveQuestionType(type);
    setCurrentPage(1);
  };

  const handleSearch = (e) => {
    setSearchQuery(e.target.value);
    console.log(searchQuery);
  };

  const handleSearchSubmit = () => {
    // 点击搜索时重置到第一页
    setCurrentPage(1);
  };

  const toggleAddQuestionDropdown = () => {
    setShowAddQuestionDropdown(!showAddQuestionDropdown);
  };

  const toggleQuestionSelection = (id) => {
    setQuestions(questions.map(q =>
        q.id === id ? { ...q, selected: !q.selected } : q
    ));

    setSelectedQuestions(prev => {
      if (prev.includes(id)) {
        return prev.filter(qId => qId !== id);
      } else {
        return [...prev, id];
      }
    });
  };
    // 这里实际应该打开对应的添加题目模态框
  const handleAddQuestion = async (type) => {
    if (type === 'manual') {
      setEditingQuestion({
        id: -Date.now(), // 临时ID
        type: '单选题',
        title: '',
        content: '',
        creator: '当前用户',
        selected: false,
        languages: ['js'],
        options: ['', '', '', ''], // 改为数组形式
        answer: [], // 改为数组形式
        difficulty: 'medium'
      });
      setShowEditModal(true);
    } else {
      try {
        setIsLoading(true);
        const generated = await api.generateQuestions({
          type: '单选题',
          count: 5,
          difficulty: 'medium',
          language: 'js'
        });
        setAiGeneratedQuestions(generated);
        setShowAIPreview(true);
      } catch (error) {
        console.error('AI生成题目失败:', error);
        alert('AI生成题目失败，请稍后重试');
      } finally {
        setIsLoading(false);
      }
    }
    setShowAddQuestionDropdown(false);
  };

  const handleSaveAIGeneratedQuestions = async (selectedIds) => {
    try {
      setIsLoading(true);
      const toSave = aiGeneratedQuestions.filter(q => selectedIds.includes(q.id));
      const result = await api.saveAIGeneratedQuestions(toSave);

      if (result.success) {
        const data = await api.getQuestions({
          page: currentPage,
          pageSize: itemsPerPage,
          type: activeQuestionType === '全部' ? '0' : REVERSE_QUESTION_TYPE_MAP[activeQuestionType],
          keyword: searchQuery
        });
        setQuestions(data.questions);
        setTotalQuestions(data.total);
        setShowAIPreview(false);
        alert(`成功保存${result.savedCount}道题目`);
      }
    } catch (error) {
      console.error('保存AI题目失败:', error);
      alert('保存题目失败，请稍后重试');
    } finally {
      setIsLoading(false);
    }
  };

  const saveEditedQuestion = async () => {
    if (!editingQuestion.title.trim()) {
      alert('请输入题目标题');
      return;
    }

    if (!editingQuestion.content.trim()) {
      alert('请输入题目内容');
      return;
    }

    if (editingQuestion.type !== '编程题') {
      // 验证选项
      if (editingQuestion.options.some(opt => !opt.trim())) {
        alert('所有选项都不能为空');
        return;
      }

      // 验证答案
      if (editingQuestion.answer.length === 0) {
        alert('请设置正确答案');
        return;
      }

      if (editingQuestion.type === '单选题' && editingQuestion.answer.length !== 1) {
        alert('单选题的答案有且仅有一个');
        return;
      }

      if (editingQuestion.type === '多选题' && editingQuestion.answer.length <= 1) {
        alert('多选题的答案必须大于一个');
        return;
      }
    }

    try {
      setIsLoading(true);

      const numericAnswer = editingQuestion.answer.map(a => parseInt(a));

      const questionData = {
        Title: editingQuestion.title,
        Type: REVERSE_QUESTION_TYPE_MAP[editingQuestion.type],
        Content: editingQuestion.content,
        Language: editingQuestion.languages?.[0] || 'js',
        Options: JSON.stringify(editingQuestion.options),
        Answer: JSON.stringify(editingQuestion.type === '单选题' ?
            [numericAnswer[0]] :
            numericAnswer.sort())
      };

      let result;
      if (editingQuestion.id > 0) {
        // 更新题目
        result = await api.updateQuestion(editingQuestion.id, {
          ...questionData,
          ID: editingQuestion.id
        });
      } else {
        // 创建题目
        result = await api.createQuestion(questionData);
      }

      // 刷新列表
      const data = await api.getQuestions({
        pageNum: currentPage,
        pageSize: itemsPerPage,
        type: activeQuestionType === '全部' ? '0' : REVERSE_QUESTION_TYPE_MAP[activeQuestionType],
        keyword: searchQuery
      });

      const formattedQuestions = data.data.records.map(q => ({
        ...q,
        type: QUESTION_TYPE_MAP[q.type],
        title: q.title,
        content: q.content,
        options: JSON.parse(q.options),
        id: q.id,
        answer: JSON.parse(q.answer),
        selected: false,
      }));

      setQuestions(formattedQuestions);
      setTotalQuestions(data.total);
      setShowEditModal(false);
      setEditingQuestion(null);
    } catch (error) {
      console.error('保存题目失败:', error);
      alert('保存题目失败: ' + error.message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleBatchDelete = () => {
    if (selectedQuestions.length === 0) {
      alert('请至少选择一道题目');
      return;
    }
    setShowDeleteConfirm(true);
  };

  const confirmSingleDelete = async () => {
    try {
      setIsLoading(true);
      await api.deleteQuestions(questionToDelete);

      // 刷新列表
      const data = await api.getQuestions({
        pageNum: currentPage,
        pageSize: itemsPerPage,
        type: activeQuestionType === '全部' ? '0' : REVERSE_QUESTION_TYPE_MAP[activeQuestionType],
        keyword: searchQuery
      });

      const formattedQuestions = data.data.records.map(q => ({
        ...q,
        type: QUESTION_TYPE_MAP[q.type],
        title: q.title,
        content: q.content,
        options: JSON.parse(q.options),
        id: q.id,
        answer: JSON.parse(q.answer),
        selected: false,
      }));

      setQuestions(formattedQuestions);
      setTotalQuestions(data.total);
      setShowDeleteConfirm(false);
      setQuestionToDelete(null);

      // 从已选列表中移除（如果存在）
      setSelectedQuestions(prev => prev.filter(id => id !== questionToDelete));
    } catch (error) {
      console.error('删除题目失败:', error);
      alert(`删除题目失败: ${error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const confirmDelete = async () => {
    try {
      setIsLoading(true);
      await api.deleteQuestions(selectedQuestions);

      // 刷新列表
      const data = await api.getQuestions({
        pageNum: 1, // 删除后回到第一页
        pageSize: itemsPerPage,
        type: activeQuestionType === '全部' ? '0' : REVERSE_QUESTION_TYPE_MAP[activeQuestionType],
        keyword: searchQuery
      });

      const formattedQuestions = data.data.records.map(q => ({
        ...q,
        type: QUESTION_TYPE_MAP[q.type],
        title: q.title,
        content: q.content,
        options: JSON.parse(q.options),
        id: q.id,
        answer: JSON.parse(q.answer),
        selected: false,
      }));

      setQuestions(formattedQuestions);
      setTotalQuestions(data.total);
      setSelectedQuestions([]);
      setShowDeleteConfirm(false);
    } catch (error) {
      console.error('批量删除失败:', error);
      alert(`删除题目失败: ${error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeleteQuestion = (id) => {
    // 删除题目函数
    setQuestionToDelete(id);
    setShowDeleteConfirm(true);
  };


  const handleEditQuestion = (question) => {
    //编辑题目函数
    setEditingQuestion(question);
    setShowEditModal(true);
  };



  const handlePageChange = (page) => {
    setCurrentPage(page);
  };

  const handleItemsPerPageChange = (value) => {
    setItemsPerPage(Number(value));
    setCurrentPage(1);
  };

  const filteredQuestions = questions.filter(q => {
    if (!q) return false; // 防止空对象

    const typeMatch = activeQuestionType === '全部' || q.type === activeQuestionType;

    // 安全访问 content 和 creator
    const content = q.content || '';
    const creator = q.creator || '';

    const searchMatch =
        content.toLowerCase().includes(searchQuery.toLowerCase()) ||
        creator.toLowerCase().includes(searchQuery.toLowerCase());

    return typeMatch && searchMatch;
  });

  const pageCount = Math.ceil(filteredQuestions.length / itemsPerPage);
  const displayedQuestions = filteredQuestions.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );


  return (
      <div className="app">
        {/* 头部 */}
        <header className="header">
          <h1>题库管理系统</h1>
          <button className="nav-toggle" onClick={toggleNav}>
            {isNavCollapsed ? '展开导航' : '折叠导航'}
          </button>
        </header>

        <div className="main-content">
          {/* 导航栏 */}
          <nav className={`sidebar ${isNavCollapsed ? 'collapsed' : ''}`}>
            <ul>
              <li><a href="#">首页</a></li>
              <li><a href="#">题库管理</a></li>
              <li><a href="#">用户管理</a></li>
              <li><a href="#">系统设置</a></li>
            </ul>
          </nav>

          {/* 主体内容 */}
          <main className="content">
            {/* 第一部分：题型筛选与搜索框 */}
            <div className="filter-section">
              <div className="question-type-filter">
                {['全部', '单选题', '多选题', '编程题'].map(type => (
                    <button
                        key={type}
                        className={`type-btn ${activeQuestionType === type ? 'active' : ''}`}
                        onClick={() => handleQuestionTypeClick(type)}
                    >
                      {type}
                    </button>
                ))}
              </div>
              <div className="search-box">
                <input
                    type="text"
                    placeholder="搜索题目或创建人..."
                    value={searchQuery}
                    onChange={handleSearch}
                />
                <button className="search-btn" onClick={handleSearchSubmit}>搜索</button>
              </div>
            </div>

            {/* 第二部分：出题模块 */}
            <div className="action-section">
              <div className="add-question">
                <button className="add-btn" onClick={toggleAddQuestionDropdown}>
                  出题
                </button>
                {showAddQuestionDropdown && (
                    <div className="add-question-dropdown">
                      <button onClick={() => handleAddQuestion('ai')}>AI出题</button>
                      <button onClick={() => handleAddQuestion('manual')}>自主出题</button>
                    </div>
                )}
              </div>
              <button
                  className="batch-delete-btn"
                  onClick={handleBatchDelete}
                  disabled={selectedQuestions.length === 0}
              >
                批量删除
              </button>
            </div>

            {/* 第三部分：题目显示 */}
            <div className="questions-table">
              <table>
                <thead>
                <tr>
                  <th width="50px"><input type="checkbox" /></th>
                  <th>题目标题</th>
                  <th width="100px">题型</th>
                  <th width="100px">语言</th>
                  <th width="150px">操作</th>
                </tr>
                </thead>
                <tbody>
                {displayedQuestions.map(question => (
                    <tr key={question.id}>
                      <td>
                        <input
                            type="checkbox"
                            checked={question.selected}
                            onChange={() => toggleQuestionSelection(question.id)}
                        />
                      </td>
                      <td>{question.title}</td>
                      <td>{question.type}</td>
                      <td>{question.language}</td>
                      <td>
                        <button
                            className="edit-btn"
                            onClick={() => handleEditQuestion(question)}
                        >
                          编辑
                        </button>
                        <button
                            className="delete-btn"
                            onClick={() => handleDeleteQuestion(question.id)}
                        >
                          删除
                        </button>
                      </td>
                    </tr>
                ))}
                </tbody>
              </table>
            </div>

            {/* 分页功能 */}
            <div className="pagination">
              <div className="pagination-controls">
                <button
                    disabled={currentPage === 1}
                    onClick={() => handlePageChange(currentPage - 1)}
                >
                  上一页
                </button>

                {Array.from({ length: Math.min(5, pageCount) }, (_, i) => {
                  let pageNum;
                  if (pageCount <= 5) {
                    pageNum = i + 1;
                  } else if (currentPage <= 3) {
                    pageNum = i + 1;
                  } else if (currentPage >= pageCount - 2) {
                    pageNum = pageCount - 4 + i;
                  } else {
                    pageNum = currentPage - 2 + i;
                  }

                  return (
                      <button
                          key={pageNum}
                          className={currentPage === pageNum ? 'active' : ''}
                          onClick={() => handlePageChange(pageNum)}
                      >
                        {pageNum}
                      </button>
                  );
                })}

                <button
                    disabled={currentPage === pageCount}
                    onClick={() => handlePageChange(currentPage + 1)}
                >
                  下一页
                </button>
              </div>

              <div className="page-size-selector">
                <span>每页显示：</span>
                <select
                    value={itemsPerPage}
                    onChange={(e) => handleItemsPerPageChange(e.target.value)}
                >
                  <option value="10">10条</option>
                  <option value="20">20条</option>
                  <option value="50">50条</option>
                </select>
              </div>

              <div className="page-jump">
                <span>跳转到：</span>
                <input
                    type="number"
                    min="1"
                    max={pageCount}
                    onBlur={(e) => {
                      const page = Math.min(Math.max(1, parseInt(e.target.value)), pageCount);
                      handlePageChange(page);
                      e.target.value = '';
                    }}
                />
                <span>页</span>
              </div>
            </div>
          </main>
        </div>

        {/* 删除确认弹窗 */}
        {showDeleteConfirm && (
            <div className="modal-overlay">
              <div className="confirm-modal">
                <h3>确认删除</h3>
                <p>{questionToDelete ? '确定要删除这道题目吗？' : `确定要删除选中的${selectedQuestions.length}道题目吗？`}</p>
                <div className="modal-actions">
                  <button onClick={questionToDelete ? confirmSingleDelete : confirmDelete}>确定</button>
                  <button onClick={() => {
                    setShowDeleteConfirm(false);
                    setQuestionToDelete(null);
                  }}>取消</button>
                </div>
              </div>
            </div>
        )}

        {showEditModal && (
            <div className="modal-overlay">
              <div className="edit-modal">
                <h3>{editingQuestion.id > 0 ? '编辑题目' : '新增题目'}</h3>
                <div className="form-group">
                  <label>题型：</label>
                  <select
                      value={editingQuestion.type}
                      onChange={(e) => setEditingQuestion({
                        ...editingQuestion,
                        type: e.target.value,
                        answer: e.target.value === '编程题' ? [] : editingQuestion.answer
                      })}
                  >
                    <option value="单选题">单选题</option>
                    <option value="多选题">多选题</option>
                    <option value="编程题">编程题</option>
                  </select>
                </div>

                <div className="form-group">
                  <label>编程语言：</label>
                  <select
                      value={editingQuestion.languages?.[0] || 'js'}
                      onChange={(e) => setEditingQuestion({
                        ...editingQuestion,
                        languages: [e.target.value]
                      })}
                  >
                    <option value="js">JavaScript</option>
                    <option value="go">Go</option>
                  </select>
                </div>

                <div className="form-group">
                  <label>题目标题：</label>
                  <input
                      type="text"
                      value={editingQuestion.title}
                      onChange={(e) => setEditingQuestion({
                        ...editingQuestion,
                        title: e.target.value
                      })}
                  />
                </div>

                <div className="form-group">
                  <label>题目内容：</label>
                  <textarea
                      value={editingQuestion.content}
                      onChange={(e) => setEditingQuestion({
                        ...editingQuestion,
                        content: e.target.value
                      })}
                  />
                </div>

                {editingQuestion.type !== '编程题' && (
                    <div className="question-options">
                      <h4>选项设置</h4>
                      {[0, 1, 2, 3].map(index => (
                          <div key={index} className="form-group option-item">
                            <label>选项{String.fromCharCode(65 + index)}：</label>
                            <input
                                type="text"
                                value={editingQuestion.options[index] || ''}
                                onChange={(e) => {
                                  const newOptions = [...editingQuestion.options];
                                  newOptions[index] = e.target.value;
                                  setEditingQuestion({
                                    ...editingQuestion,
                                    options: newOptions
                                  });
                                }}
                            />
                          </div>
                      ))}

                      <div className="form-group">
                        <label>正确答案：</label>
                        {editingQuestion.type === '单选题' ? (
                            <select
                                value={editingQuestion.answer[0] || ''}
                                onChange={(e) => setEditingQuestion({
                                  ...editingQuestion,
                                  answer: e.target.value ? [e.target.value] : []
                                })}
                            >
                              <option value="">请选择</option>
                              {[0, 1, 2, 3].map(index => (
                                  <option
                                      key={index}
                                      value={index}
                                      disabled={!editingQuestion.options[index]}
                                  >
                                    {String.fromCharCode(65 + index)} {!editingQuestion.options[index] && '(未填写)'}
                                  </option>
                              ))}
                            </select>
                        ) : (
                            <div className="multi-answer">
                              {[0, 1, 2, 3].map(index => (
                                  <label key={index} className="checkbox-option">
                                    <input
                                        type="checkbox"
                                        checked={editingQuestion.answer.includes(index)}
                                        onChange={(e) => {
                                          let newAnswer;
                                          if (e.target.checked) {
                                            newAnswer = [...editingQuestion.answer, index].sort();
                                          } else {
                                            newAnswer = editingQuestion.answer.filter(a => a !== index);
                                          }
                                          setEditingQuestion({
                                            ...editingQuestion,
                                            answer: newAnswer
                                          });
                                        }}
                                        disabled={!editingQuestion.options[index]}
                                    />
                                    <span>
                            {String.fromCharCode(65 + index)} {!editingQuestion.options[index] && '(未填写)'}
                          </span>
                                  </label>
                              ))}
                            </div>
                        )}
                      </div>
                    </div>
                )}

                <div className="modal-actions">
                  <button onClick={saveEditedQuestion}>保存</button>
                  <button onClick={() => setShowEditModal(false)}>取消</button>
                </div>
              </div>
            </div>
        )}
      </div>
  );
}

export default App;