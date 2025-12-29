(function() {
  'use strict';

  // Get chatbot ID from script tag
  // Look for the script tag that loaded this file
  let chatbotId = null;
  let apiUrl = null;
  
  const scripts = document.getElementsByTagName('script');
  for (let i = 0; i < scripts.length; i++) {
    const script = scripts[i];
    if (script.src && script.src.includes('widget.js')) {
      chatbotId = script.getAttribute('data-chatbot-id');
      apiUrl = script.getAttribute('data-api-url');
      break;
    }
  }
  
  // Fallback: check for inline configuration
  if (!chatbotId && window.ChatbotConfig) {
    chatbotId = window.ChatbotConfig.chatbotId;
    apiUrl = window.ChatbotConfig.apiUrl;
  }
  
  // Set default API URL if not provided
  if (!apiUrl) {
    // Auto-detect API URL based on where the widget was loaded from
    const currentDomain = window.location.hostname;
    if (currentDomain === 'localhost' || currentDomain === '127.0.0.1') {
      apiUrl = 'http://localhost:8081/api';
    } else {
      // Use the same domain as the current page with https
      apiUrl = `${window.location.protocol}//chatbot-api.appster.co.in/api`;
    }
  }

  if (!chatbotId) {
    console.error('Chatbot ID not provided. Please add data-chatbot-id attribute to the script tag.');
    return;
  }
  
  console.log('Chatbot initialized with ID:', chatbotId, 'API URL:', apiUrl);

  // Generate unique session ID
  function generateSessionId() {
    let sessionId = localStorage.getItem('chatbot_session_' + chatbotId);
    if (!sessionId) {
      sessionId = 'sess_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
      localStorage.setItem('chatbot_session_' + chatbotId, sessionId);
    }
    return sessionId;
  }

  const sessionId = generateSessionId();
  let settings = null;
  let isOpen = false;
  let messages = [];

  // Fetch chatbot settings
  async function fetchSettings() {
    try {
      const response = await fetch(`${apiUrl}/chatbots/${chatbotId}/settings`);
      if (response.ok) {
        settings = await response.json();
      } else {
        settings = getDefaultSettings();
      }
    } catch (error) {
      console.error('Failed to fetch chatbot settings:', error);
      settings = getDefaultSettings();
    }
  }

  function getDefaultSettings() {
    return {
      theme_color: '#007bff',
      position: 'bottom-right',
      welcome_message: 'Hi! How can I help you today?',
      avatar_url: '',
      widget_size: 'medium'
    };
  }

  // Send message to backend
  async function sendMessage(message) {
    try {
      const response = await fetch(`${apiUrl}/chat/${chatbotId}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          session_id: sessionId,
          message: message
        })
      });

      if (response.ok) {
        const data = await response.json();
        return data.response;
      } else {
        return 'Sorry, I encountered an error. Please try again.';
      }
    } catch (error) {
      console.error('Failed to send message:', error);
      return 'Sorry, I encountered an error. Please try again.';
    }
  }

  // Create widget UI
  function createWidget() {
    // Create widget container
    const container = document.createElement('div');
    container.id = 'chatbot-widget-container';
    container.style.cssText = `
      position: fixed;
      ${settings.position === 'bottom-left' ? 'left' : 'right'}: 20px;
      bottom: 20px;
      z-index: 9999;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    `;

    // Create chat button
    const button = document.createElement('button');
    button.id = 'chatbot-toggle-button';
    button.innerHTML = '💬';
    button.style.cssText = `
      width: 60px;
      height: 60px;
      border-radius: 50%;
      background-color: ${settings.theme_color};
      color: white;
      border: none;
      font-size: 24px;
      cursor: pointer;
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
      transition: transform 0.2s;
    `;

    button.onmouseover = () => {
      button.style.transform = 'scale(1.1)';
    };
    button.onmouseout = () => {
      button.style.transform = 'scale(1)';
    };
    button.onclick = toggleChat;

    // Create chat window
    const chatWindow = document.createElement('div');
    chatWindow.id = 'chatbot-chat-window';
    chatWindow.style.cssText = `
      display: none;
      width: ${settings.widget_size === 'small' ? '300px' : settings.widget_size === 'large' ? '400px' : '350px'};
      height: ${settings.widget_size === 'small' ? '400px' : settings.widget_size === 'large' ? '600px' : '500px'};
      background: white;
      border-radius: 10px;
      box-shadow: 0 5px 40px rgba(0,0,0,0.16);
      display: flex;
      flex-direction: column;
      overflow: hidden;
      margin-bottom: 10px;
    `;

    // Chat header
    const header = document.createElement('div');
    header.style.cssText = `
      background-color: ${settings.theme_color};
      color: white;
      padding: 15px;
      font-weight: bold;
      display: flex;
      justify-content: space-between;
      align-items: center;
    `;
    header.innerHTML = `
      <span>Chat with us</span>
      <button id="chatbot-close-button" style="background: none; border: none; color: white; font-size: 20px; cursor: pointer;">×</button>
    `;

    // Messages container
    const messagesContainer = document.createElement('div');
    messagesContainer.id = 'chatbot-messages';
    messagesContainer.style.cssText = `
      flex: 1;
      overflow-y: auto;
      padding: 15px;
      background-color: #f5f5f5;
    `;

    // Input container
    const inputContainer = document.createElement('div');
    inputContainer.style.cssText = `
      display: flex;
      padding: 15px;
      background: white;
      border-top: 1px solid #e0e0e0;
    `;

    const input = document.createElement('input');
    input.id = 'chatbot-input';
    input.type = 'text';
    input.placeholder = 'Type your message...';
    input.style.cssText = `
      flex: 1;
      padding: 10px;
      border: 1px solid #ddd;
      border-radius: 20px;
      outline: none;
      font-size: 14px;
    `;

    const sendButton = document.createElement('button');
    sendButton.id = 'chatbot-send-button';
    sendButton.innerHTML = '➤';
    sendButton.style.cssText = `
      background-color: ${settings.theme_color};
      color: white;
      border: none;
      border-radius: 50%;
      width: 40px;
      height: 40px;
      margin-left: 10px;
      cursor: pointer;
      font-size: 18px;
      transition: opacity 0.3s;
    `;
    sendButton.onclick = handleSendMessage;

    inputContainer.appendChild(input);
    inputContainer.appendChild(sendButton);

    chatWindow.appendChild(header);
    chatWindow.appendChild(messagesContainer);
    chatWindow.appendChild(inputContainer);

    container.appendChild(chatWindow);
    container.appendChild(button);

    document.body.appendChild(container);

    // Event listeners
    document.getElementById('chatbot-close-button').onclick = toggleChat;
    input.addEventListener('keypress', (e) => {
      if (e.key === 'Enter') {
        handleSendMessage();
      }
    });

    // Add welcome message
    addMessage('assistant', settings.welcome_message);
  }

  function toggleChat() {
    isOpen = !isOpen;
    const chatWindow = document.getElementById('chatbot-chat-window');
    const button = document.getElementById('chatbot-toggle-button');

    if (isOpen) {
      chatWindow.style.display = 'flex';
      button.style.display = 'none';
    } else {
      chatWindow.style.display = 'none';
      button.style.display = 'block';
    }
  }

  function addMessage(role, content) {
    messages.push({ role, content });
    const messagesContainer = document.getElementById('chatbot-messages');

    const messageDiv = document.createElement('div');
    messageDiv.style.cssText = `
      margin-bottom: 12px;
      display: flex;
      align-items: flex-end;
      gap: 8px;
      ${role === 'user' ? 'justify-content: flex-end;' : 'justify-content: flex-start;'}
    `;

    // Add avatar for bot messages if available
    if (role === 'assistant' && settings.avatar_url) {
      const avatar = document.createElement('img');
      avatar.src = settings.avatar_url;
      avatar.style.cssText = `
        width: 32px;
        height: 32px;
        border-radius: 50%;
        object-fit: cover;
        flex-shrink: 0;
      `;
      avatar.onerror = function() {
        this.style.display = 'none';
      };
      messageDiv.appendChild(avatar);
    }

    const messageBubble = document.createElement('div');
    messageBubble.style.cssText = `
      max-width: 70%;
      padding: 10px 15px;
      border-radius: 18px;
      ${role === 'user' 
        ? `background-color: ${settings.theme_color}; color: white; font-weight: 500;` 
        : 'background-color: #ffffff; color: #333; border: 1px solid #e0e0e0;'}
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
      font-size: 14px;
      line-height: 1.5;
      word-wrap: break-word;
    `;
    
    // For bot messages, convert URLs to clickable links
    if (role === 'assistant') {
      const formattedContent = formatMessageContent(content);
      messageBubble.innerHTML = formattedContent;
    } else {
      messageBubble.textContent = content;
    }

    messageDiv.appendChild(messageBubble);
    messagesContainer.appendChild(messageDiv);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  }

  function formatMessageContent(content) {
    // Escape HTML
    const escapeHtml = (text) => {
      const div = document.createElement('div');
      div.textContent = text;
      return div.innerHTML;
    };

    // Convert line breaks to <br>
    let formatted = escapeHtml(content).replace(/\n/g, '<br>');

    // Convert URLs to clickable links
    const urlRegex = /(https?:\/\/[^\s]+)/g;
    formatted = formatted.replace(urlRegex, (url) => {
      return `<a href="${url}" target="_blank" rel="noopener noreferrer" style="color: #1976d2; text-decoration: underline; word-break: break-all;">${url}</a>`;
    });

    return formatted;
  }

  function addTypingIndicator() {
    const messagesContainer = document.getElementById('chatbot-messages');
    const typingDiv = document.createElement('div');
    typingDiv.id = 'chatbot-typing';
    typingDiv.style.cssText = `
      margin-bottom: 12px;
      display: flex;
      align-items: flex-end;
      gap: 8px;
      justify-content: flex-start;
    `;

    // Add avatar if available
    if (settings.avatar_url) {
      const avatar = document.createElement('img');
      avatar.src = settings.avatar_url;
      avatar.style.cssText = `
        width: 32px;
        height: 32px;
        border-radius: 50%;
        object-fit: cover;
        flex-shrink: 0;
      `;
      avatar.onerror = function() {
        this.style.display = 'none';
      };
      typingDiv.appendChild(avatar);
    }

    const typingBubble = document.createElement('div');
    typingBubble.style.cssText = `
      background-color: #ffffff;
      padding: 10px 15px;
      border-radius: 18px;
      border: 1px solid #e0e0e0;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
      display: flex;
      align-items: center;
      gap: 4px;
    `;
    typingBubble.innerHTML = `
      <span style="color: #666; font-size: 14px; margin-right: 4px;">Typing</span>
      <span style="animation: blink 1.4s infinite; color: ${settings.theme_color}; font-size: 20px; line-height: 1;">●</span>
      <span style="animation: blink 1.4s infinite 0.2s; color: ${settings.theme_color}; font-size: 20px; line-height: 1;">●</span>
      <span style="animation: blink 1.4s infinite 0.4s; color: ${settings.theme_color}; font-size: 20px; line-height: 1;">●</span>
    `;

    typingDiv.appendChild(typingBubble);
    messagesContainer.appendChild(typingDiv);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  }

  function removeTypingIndicator() {
    const typing = document.getElementById('chatbot-typing');
    if (typing) {
      typing.remove();
    }
  }

  async function handleSendMessage() {
    const input = document.getElementById('chatbot-input');
    const sendButton = document.getElementById('chatbot-send-button');
    const message = input.value.trim();

    if (!message) return;

    // Disable input and button
    input.disabled = true;
    sendButton.disabled = true;
    input.style.opacity = '0.6';
    input.style.cursor = 'not-allowed';
    sendButton.style.opacity = '0.6';
    sendButton.style.cursor = 'not-allowed';

    // Add user message
    addMessage('user', message);
    input.value = '';

    // Show typing indicator
    addTypingIndicator();

    // Send to backend
    const response = await sendMessage(message);

    // Remove typing indicator
    removeTypingIndicator();

    // Add bot response
    addMessage('assistant', response);

    // Re-enable input and button
    input.disabled = false;
    sendButton.disabled = false;
    input.style.opacity = '1';
    input.style.cursor = 'text';
    sendButton.style.opacity = '1';
    sendButton.style.cursor = 'pointer';
    input.focus();
  }

  // Add animation styles
  const style = document.createElement('style');
  style.textContent = `
    @keyframes blink {
      0%, 50%, 100% { opacity: 0.3; }
      25% { opacity: 1; }
    }
  `;
  document.head.appendChild(style);

  // Initialize widget
  fetchSettings().then(() => {
    createWidget();
  });
})();

