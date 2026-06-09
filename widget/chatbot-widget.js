(function() {
  'use strict';

  // Get chatbot ID from script tag
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

  if (!chatbotId && window.ChatbotConfig) {
    chatbotId = window.ChatbotConfig.chatbotId;
    apiUrl = window.ChatbotConfig.apiUrl;
  }

  if (!apiUrl) {
    const currentDomain = window.location.hostname;
    if (currentDomain === 'localhost' || currentDomain === '127.0.0.1') {
      apiUrl = 'http://localhost:8081/api';
    } else {
      apiUrl = `${window.location.protocol}//chatbot-api.appster.co.in/api`;
    }
  }

  if (!chatbotId) {
    console.error('Chatbot ID not provided. Please add data-chatbot-id attribute to the script tag.');
    return;
  }

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
  // true once the lead form has been submitted (or is not needed)
  let leadCaptured = false;

  function hasSubmittedLead() {
    return localStorage.getItem('chatbot_lead_' + chatbotId) === '1';
  }

  function markLeadSubmitted() {
    localStorage.setItem('chatbot_lead_' + chatbotId, '1');
  }

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
      widget_size: 'medium',
      lead_capture: { enabled: false, title: '', subtitle: '', fields: [] }
    };
  }

  async function submitLead(fieldValues) {
    try {
      const response = await fetch(`${apiUrl}/leads/${chatbotId}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ session_id: sessionId, field_values: fieldValues })
      });
      return response.ok;
    } catch (e) {
      console.error('Failed to submit lead:', e);
      return false;
    }
  }

  async function sendMessage(message) {
    try {
      const response = await fetch(`${apiUrl}/chat/${chatbotId}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ session_id: sessionId, message: message })
      });
      if (response.ok) {
        const data = await response.json();
        return data.response;
      }
      return 'Sorry, I encountered an error. Please try again.';
    } catch (error) {
      console.error('Failed to send message:', error);
      return 'Sorry, I encountered an error. Please try again.';
    }
  }

  // ─── Widget creation ────────────────────────────────────────────────────────

  function createWidget() {
    // Inject theme-coloured animations now that settings are loaded
    const tc = settings.theme_color || '#007bff';
    style.textContent += `
      @keyframes chatbot-msg-in {
        0%   { opacity:0; transform:translateX(-10px); box-shadow:0 0 0 2.5px ${tc}88, 0 4px 14px ${tc}55; }
        35%  { opacity:1; transform:translateX(0);     box-shadow:0 0 0 2.5px ${tc}88, 0 4px 14px ${tc}55; }
        70%  { box-shadow:0 0 0 2px ${tc}44, 0 2px 8px ${tc}33; }
        100% { opacity:1; transform:translateX(0);     box-shadow:0 2px 4px rgba(0,0,0,0.1); }
      }
      .chatbot-response-bubble {
        animation: chatbot-msg-in 1.6s cubic-bezier(0.22,1,0.36,1) forwards;
      }
    `;

    const container = document.createElement('div');
    container.id = 'chatbot-widget-container';
    container.style.cssText = `
      position: fixed;
      ${settings.position === 'bottom-left' ? 'left' : 'right'}: 20px;
      bottom: 20px;
      z-index: 9999;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    `;

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
    button.onmouseover = () => { button.style.transform = 'scale(1.1)'; };
    button.onmouseout  = () => { button.style.transform = 'scale(1)'; };
    button.onclick = toggleChat;

    const chatWindow = document.createElement('div');
    chatWindow.id = 'chatbot-chat-window';
    const w = settings.widget_size === 'small' ? '300px' : settings.widget_size === 'large' ? '400px' : '350px';
    const h = settings.widget_size === 'small' ? '400px' : settings.widget_size === 'large' ? '600px' : '500px';
    chatWindow.style.cssText = `
      display: none;
      position: absolute;
      bottom: 80px;
      ${settings.position === 'bottom-left' ? 'left' : 'right'}: 0;
      width: ${w};
      height: ${h};
      background: white;
      border-radius: 10px;
      box-shadow: 0 5px 40px rgba(0,0,0,0.16);
      flex-direction: column;
      overflow: hidden;
    `;

    // Header
    const header = document.createElement('div');
    header.style.cssText = `
      background-color: ${settings.theme_color};
      color: white;
      padding: 15px;
      font-weight: bold;
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-shrink: 0;
    `;
    header.innerHTML = `
      <span>Chat with us</span>
      <button id="chatbot-close-button" style="background:none;border:none;color:white;font-size:20px;cursor:pointer;">×</button>
    `;

    chatWindow.appendChild(header);

    const lc = settings.lead_capture || {};
    const needsLeadForm = lc.enabled && !hasSubmittedLead() && lc.fields && lc.fields.length > 0;

    if (needsLeadForm) {
      chatWindow.appendChild(buildLeadForm(lc, chatWindow, header));
    } else {
      leadCaptured = true;
      chatWindow.appendChild(buildChatBody());
    }

    container.appendChild(chatWindow);
    container.appendChild(button);
    document.body.appendChild(container);

    if (!needsLeadForm) {
      document.getElementById('chatbot-close-button').onclick = toggleChat;
    }
  }

  // ─── Lead capture form ──────────────────────────────────────────────────────

  function buildLeadForm(lc, chatWindow, defaultHeader) {
    // Hide the plain default header — our custom header replaces it
    defaultHeader.style.display = 'none';

    const wrapper = document.createElement('div');
    wrapper.id = 'chatbot-lead-form-wrapper';
    wrapper.style.cssText = 'display:flex;flex-direction:column;flex:1;overflow:hidden;';

    // ── Premium dark header ──────────────────────────────────────────────────
    const hdr = document.createElement('div');
    hdr.style.cssText = `
      background: linear-gradient(150deg, #1a2b52 0%, #243368 55%, #2d3f80 100%);
      padding: 18px 18px 22px;
      position: relative;
      flex-shrink: 0;
      border-radius: 0 0 22px 22px;
    `;

    // Sparkle SVG icon
    const sparkle = document.createElement('div');
    sparkle.innerHTML = `<svg width="30" height="30" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M12 2L13.09 9.26L20 12L13.09 14.74L12 22L10.91 14.74L4 12L10.91 9.26L12 2Z" fill="rgba(255,255,255,0.95)"/>
      <path d="M19 5L19.54 8.46L23 9L19.54 9.54L19 13L18.46 9.54L15 9L18.46 8.46L19 5Z" fill="rgba(255,255,255,0.6)"/>
    </svg>`;
    sparkle.style.cssText = 'margin-bottom:12px;';

    const closeBtn = document.createElement('button');
    closeBtn.innerHTML = '&times;';
    closeBtn.style.cssText = `
      position:absolute;top:10px;right:13px;
      background:rgba(255,255,255,0.12);border:none;color:white;
      font-size:20px;cursor:pointer;line-height:1;padding:3px 7px;
      border-radius:6px;transition:background 0.2s;
    `;
    closeBtn.onmouseover = () => { closeBtn.style.background = 'rgba(255,255,255,0.22)'; };
    closeBtn.onmouseout  = () => { closeBtn.style.background = 'rgba(255,255,255,0.12)'; };
    closeBtn.onclick = toggleChat;

    const askLabel = document.createElement('div');
    askLabel.textContent = 'Ask a question';
    askLabel.style.cssText = 'color:rgba(255,255,255,0.65);font-size:12px;font-weight:500;letter-spacing:0.5px;margin-bottom:5px;';

    const mainTitle = document.createElement('div');
    mainTitle.textContent = lc.title || 'How can we help?';
    mainTitle.style.cssText = 'color:#ffffff;font-size:20px;font-weight:700;line-height:1.25;';

    hdr.appendChild(sparkle);
    hdr.appendChild(closeBtn);
    hdr.appendChild(askLabel);
    hdr.appendChild(mainTitle);
    wrapper.appendChild(hdr);

    // ── Scrollable body ──────────────────────────────────────────────────────
    const body = document.createElement('div');
    body.style.cssText = `
      flex:1;min-height:0;overflow-y:auto;padding:14px;
      background:#f4f6fb;
      display:flex;flex-direction:column;gap:12px;
    `;

    // Info card (subtitle)
    const card = document.createElement('div');
    card.style.cssText = `
      background:white;
      border:1.5px solid ${settings.theme_color};
      border-radius:14px;
      padding:16px 14px;
      text-align:center;
      box-shadow:0 2px 10px rgba(0,0,0,0.06);
    `;
    const cardIcon = document.createElement('div');
    cardIcon.innerHTML = `<svg width="22" height="22" viewBox="0 0 24 24" fill="none">
      <path d="M12 2L13.09 9.26L20 12L13.09 14.74L12 22L10.91 14.74L4 12L10.91 9.26L12 2Z" fill="${settings.theme_color}"/>
      <path d="M19 5L19.54 8.46L23 9L19.54 9.54L19 13L18.46 9.54L15 9L18.46 8.46L19 5Z" fill="${settings.theme_color}99"/>
    </svg>`;
    cardIcon.style.cssText = 'display:flex;justify-content:center;margin-bottom:9px;';

    const cardTitle = document.createElement('div');
    cardTitle.textContent = lc.subtitle || 'Please enter your details for a better experience';
    cardTitle.style.cssText = 'font-size:13px;font-weight:600;color:#1a2b52;line-height:1.45;';

    card.appendChild(cardIcon);
    card.appendChild(cardTitle);
    body.appendChild(card);

    // Form fields
    lc.fields.forEach(field => {
      const fw = document.createElement('div');
      fw.style.cssText = 'display:flex;flex-direction:column;gap:5px;';

      const lbl = document.createElement('label');
      lbl.style.cssText = 'font-size:13px;font-weight:500;color:#444;';
      lbl.textContent = 'Enter ' + field.label + (field.required ? '' : ' (optional)');

      let inp;
      if (field.type === 'textarea') {
        inp = document.createElement('textarea');
        inp.rows = 3;
        inp.style.cssText = inputStyle() + 'resize:none;';
      } else {
        inp = document.createElement('input');
        inp.type = field.type || 'text';
        inp.style.cssText = inputStyle();
      }
      inp.name = field.name;
      inp.placeholder = field.placeholder || '';
      inp.required = !!field.required;
      inp.onfocus = () => { inp.style.borderColor = settings.theme_color; inp.style.boxShadow = `0 0 0 3px ${settings.theme_color}22`; };
      inp.onblur  = () => { inp.style.borderColor = '#dde1ea'; inp.style.boxShadow = 'none'; };

      fw.appendChild(lbl);
      fw.appendChild(inp);
      body.appendChild(fw);
    });

    // Error message
    const errEl = document.createElement('div');
    errEl.id = 'chatbot-lead-error';
    errEl.style.cssText = 'color:#d32f2f;font-size:12px;display:none;padding:4px 0;';
    body.appendChild(errEl);

    wrapper.appendChild(body);

    // ── Footer / submit button ───────────────────────────────────────────────
    const footer = document.createElement('div');
    footer.style.cssText = 'padding:10px 16px 14px;background:#f4f6fb;flex-shrink:0;';

    const submitBtn = document.createElement('button');
    submitBtn.textContent = 'Start Chat';
    submitBtn.style.cssText = `
      width:100%;padding:13px;
      background:linear-gradient(135deg,#1a2b52 0%,#2d3f80 100%);
      color:white;border:none;border-radius:10px;
      font-size:15px;font-weight:600;cursor:pointer;
      letter-spacing:0.3px;transition:opacity 0.2s;
      box-shadow:0 4px 14px rgba(26,43,82,0.35);
    `;
    submitBtn.onmouseover = () => { submitBtn.style.opacity = '0.88'; };
    submitBtn.onmouseout  = () => { submitBtn.style.opacity = '1'; };

    submitBtn.onclick = async () => {
      const fieldValues = {};
      let valid = true;
      lc.fields.forEach(field => {
        const el = body.querySelector(`[name="${field.name}"]`);
        const val = el ? el.value.trim() : '';
        if (field.required && !val) valid = false;
        fieldValues[field.name] = val;
      });

      if (!valid) {
        errEl.textContent = 'Please fill in all required fields.';
        errEl.style.display = 'block';
        return;
      }
      errEl.style.display = 'none';
      submitBtn.disabled = true;
      submitBtn.textContent = 'Please wait...';

      await submitLead(fieldValues);
      markLeadSubmitted();
      leadCaptured = true;

      // Restore default header, wire close button, swap form → chat
      defaultHeader.style.display = 'flex';
      document.getElementById('chatbot-close-button').onclick = toggleChat;
      chatWindow.replaceChild(buildChatBody(), wrapper);
    };

    // Disclaimer inside footer
    const leadDisclaimer = document.createElement('div');
    leadDisclaimer.style.cssText = 'text-align:center;padding-top:10px;font-size:11px;color:#aaa;';
    const leadBotName = settings.chatbot_name || 'AI';
    leadDisclaimer.innerHTML = `<a href="#" onclick="return false;" style="color:${settings.theme_color};text-decoration:none;font-weight:500;">${leadBotName}</a> can make mistakes. Verify important info.`;

    footer.appendChild(submitBtn);
    footer.appendChild(leadDisclaimer);
    wrapper.appendChild(footer);
    return wrapper;
  }

  function inputStyle() {
    return `
      padding: 11px 14px;
      border: 1.5px solid #dde1ea;
      border-radius: 10px;
      font-size: 14px;
      outline: none;
      font-family: inherit;
      background: white;
      width: 100%;
      box-sizing: border-box;
      transition: border-color 0.2s, box-shadow 0.2s;
    `;
  }

  // ─── Chat body ──────────────────────────────────────────────────────────────

  function buildChatBody() {
    const wrapper = document.createElement('div');
    wrapper.id = 'chatbot-body-wrapper';
    wrapper.style.cssText = 'display:flex;flex-direction:column;flex:1;overflow:hidden;';

    const messagesContainer = document.createElement('div');
    messagesContainer.id = 'chatbot-messages';
    messagesContainer.style.cssText = `
      flex: 1;
      min-height: 0;
      overflow-y: auto;
      padding: 15px;
      background-color: #f5f5f5;
    `;

    const inputContainer = document.createElement('div');
    inputContainer.style.cssText = `
      display: flex;
      flex-direction: column;
      padding: 10px 15px 8px;
      background: white;
      border-top: 1px solid #e0e0e0;
      flex-shrink: 0;
    `;

    // Input row
    const inputRow = document.createElement('div');
    inputRow.style.cssText = 'display:flex;align-items:center;gap:0;';

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
      flex-shrink: 0;
    `;
    sendButton.onclick = handleSendMessage;
    input.addEventListener('keypress', (e) => {
      if (e.key === 'Enter') handleSendMessage();
    });

    inputRow.appendChild(input);
    inputRow.appendChild(sendButton);

    // Disclaimer footer inside the input container
    const disclaimer = document.createElement('div');
    disclaimer.style.cssText = 'text-align:center;padding-top:6px;font-size:11px;color:#aaa;';
    const botName = settings.chatbot_name || 'AI';
    disclaimer.innerHTML = `<a href="#" onclick="return false;" style="color:${settings.theme_color};text-decoration:none;font-weight:500;">${botName}</a> can make mistakes. Verify important info.`;

    inputContainer.appendChild(inputRow);
    inputContainer.appendChild(disclaimer);

    wrapper.appendChild(messagesContainer);
    wrapper.appendChild(inputContainer);

    // Show welcome message after DOM is inserted (next tick)
    setTimeout(() => addMessage('assistant', settings.welcome_message), 0);

    return wrapper;
  }

  // ─── Toggle ─────────────────────────────────────────────────────────────────

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

  // ─── Messaging ──────────────────────────────────────────────────────────────

  function addMessage(role, content) {
    messages.push({ role, content });
    const messagesContainer = document.getElementById('chatbot-messages');
    if (!messagesContainer) return;

    const messageDiv = document.createElement('div');
    messageDiv.style.cssText = `
      margin-bottom: 12px;
      display: flex;
      align-items: flex-end;
      gap: 8px;
      ${role === 'user' ? 'justify-content: flex-end;' : 'justify-content: flex-start;'}
    `;

    if (role === 'assistant' && settings.avatar_url) {
      const avatar = document.createElement('img');
      avatar.src = settings.avatar_url;
      avatar.style.cssText = 'width:32px;height:32px;border-radius:50%;object-fit:cover;flex-shrink:0;';
      avatar.onerror = function() { this.style.display = 'none'; };
      messageDiv.appendChild(avatar);
    }

    const bubble = document.createElement('div');
    bubble.style.cssText = `
      max-width: ${role === 'user' ? '70%' : '85%'};
      padding: 10px 14px;
      border-radius: 18px;
      ${role === 'user'
        ? `background-color: ${settings.theme_color}; color: white; font-weight: 500;`
        : 'background-color: #ffffff; color: #333; border: 1px solid #e0e0e0;'}
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
      font-size: 14px;
      line-height: 1.6;
      word-wrap: break-word;
    `;

    if (role === 'assistant') {
      bubble.classList.add('chatbot-response-bubble');
      bubble.innerHTML = formatMessageContent(content);
    } else {
      bubble.textContent = content;
    }

    messageDiv.appendChild(bubble);
    messagesContainer.appendChild(messageDiv);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  }

  function formatMessageContent(content) {
    const esc = (t) => t.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');

    const inline = (t) => {
      // bold
      t = t.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
      // inline code
      t = t.replace(/`([^`]+)`/g, '<code style="background:#f0f0f0;padding:1px 5px;border-radius:3px;font-size:12px;font-family:monospace;">$1</code>');
      // links
      t = t.replace(/(https?:\/\/[^\s<]+)/g, '<a href="$1" target="_blank" rel="noopener noreferrer" style="color:#1976d2;text-decoration:underline;word-break:break-all;">$1</a>');
      return t;
    };

    const lines = content.split('\n');
    let html = '';
    let inUL = false;
    let inOL = false;

    const closeUL = () => { if (inUL) { html += '</ul>'; inUL = false; } };
    const closeOL = () => { if (inOL) { html += '</ol>'; inOL = false; } };
    const closeLists = () => { closeUL(); closeOL(); };

    lines.forEach(line => {
      const h3 = line.match(/^###\s+(.+)/);
      const h2 = line.match(/^##\s+(.+)/);
      const h1 = line.match(/^#\s+(.+)/);
      const ul = line.match(/^[-*]\s+(.+)/);
      const ol = line.match(/^\d+[.)]\s+(.+)/);

      if (h3) {
        closeLists();
        html += `<h3 style="font-size:13px;font-weight:700;margin:8px 0 3px;color:#1a1a1a;border-bottom:1px solid #eee;padding-bottom:2px;">${inline(esc(h3[1]))}</h3>`;
      } else if (h2) {
        closeLists();
        html += `<h2 style="font-size:14px;font-weight:700;margin:10px 0 4px;color:#1a1a1a;border-bottom:1px solid #eee;padding-bottom:2px;">${inline(esc(h2[1]))}</h2>`;
      } else if (h1) {
        closeLists();
        html += `<h1 style="font-size:15px;font-weight:700;margin:10px 0 5px;color:#1a1a1a;">${inline(esc(h1[1]))}</h1>`;
      } else if (ul) {
        closeOL();
        if (!inUL) { html += '<ul style="margin:4px 0;padding-left:20px;list-style-type:disc;">'; inUL = true; }
        html += `<li style="margin-bottom:4px;">${inline(esc(ul[1]))}</li>`;
      } else if (ol) {
        closeUL();
        if (!inOL) { html += '<ol style="margin:4px 0;padding-left:20px;">'; inOL = true; }
        html += `<li style="margin-bottom:4px;">${inline(esc(ol[1]))}</li>`;
      } else if (line.trim() === '') {
        closeLists();
        html += '<div style="height:5px;"></div>';
      } else {
        closeLists();
        html += `<p style="margin:0 0 4px;">${inline(esc(line))}</p>`;
      }
    });

    closeLists();
    return html;
  }

  function addTypingIndicator() {
    const messagesContainer = document.getElementById('chatbot-messages');
    if (!messagesContainer) return;
    const typingDiv = document.createElement('div');
    typingDiv.id = 'chatbot-typing';
    typingDiv.style.cssText = 'margin-bottom:12px;display:flex;align-items:flex-end;gap:8px;justify-content:flex-start;';

    if (settings.avatar_url) {
      const avatar = document.createElement('img');
      avatar.src = settings.avatar_url;
      avatar.style.cssText = 'width:32px;height:32px;border-radius:50%;object-fit:cover;flex-shrink:0;';
      avatar.onerror = function() { this.style.display = 'none'; };
      typingDiv.appendChild(avatar);
    }

    const bubble = document.createElement('div');
    bubble.style.cssText = `
      background-color: #ffffff;
      padding: 10px 15px;
      border-radius: 18px;
      border: 1px solid #e0e0e0;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
      display: flex;
      align-items: center;
      gap: 4px;
    `;
    bubble.innerHTML = `
      <span style="color:#666;font-size:14px;margin-right:4px;">Typing</span>
      <span style="animation:blink 1.4s infinite;color:${settings.theme_color};font-size:20px;line-height:1;">●</span>
      <span style="animation:blink 1.4s infinite 0.2s;color:${settings.theme_color};font-size:20px;line-height:1;">●</span>
      <span style="animation:blink 1.4s infinite 0.4s;color:${settings.theme_color};font-size:20px;line-height:1;">●</span>
    `;
    typingDiv.appendChild(bubble);
    messagesContainer.appendChild(typingDiv);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  }

  function removeTypingIndicator() {
    const el = document.getElementById('chatbot-typing');
    if (el) el.remove();
  }

  async function handleSendMessage() {
    const input = document.getElementById('chatbot-input');
    const sendButton = document.getElementById('chatbot-send-button');
    if (!input) return;
    const message = input.value.trim();
    if (!message) return;

    input.disabled = true;
    sendButton.disabled = true;
    input.style.opacity = '0.6';
    sendButton.style.opacity = '0.6';

    addMessage('user', message);
    input.value = '';
    addTypingIndicator();

    const response = await sendMessage(message);
    removeTypingIndicator();
    addMessage('assistant', response);

    input.disabled = false;
    sendButton.disabled = false;
    input.style.opacity = '1';
    sendButton.style.opacity = '1';
    input.focus();
  }

  // Styles (theme-specific keyframes are appended in createWidget once settings load)
  const style = document.createElement('style');
  style.textContent = `
    @keyframes blink {
      0%, 50%, 100% { opacity: 0.3; }
      25% { opacity: 1; }
    }
  `;
  document.head.appendChild(style);

  // Boot
  fetchSettings().then(() => {
    createWidget();
  });
})();
