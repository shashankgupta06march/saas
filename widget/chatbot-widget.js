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
    button.innerHTML = `<svg width="26" height="26" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" fill="white"/></svg>`;
    button.style.cssText = `
      width: 58px;
      height: 58px;
      border-radius: 50%;
      background: linear-gradient(135deg, ${settings.theme_color} 0%, ${settings.theme_color}cc 100%);
      color: white;
      border: none;
      cursor: pointer;
      box-shadow: 0 6px 20px ${settings.theme_color}66;
      transition: transform 0.2s, box-shadow 0.2s;
      display: flex;
      align-items: center;
      justify-content: center;
    `;
    button.onmouseover = () => { button.style.transform = 'scale(1.1)'; button.style.boxShadow = `0 8px 28px ${settings.theme_color}88`; };
    button.onmouseout  = () => { button.style.transform = 'scale(1)';   button.style.boxShadow = `0 6px 20px ${settings.theme_color}66`; };
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
      border-radius: 16px;
      box-shadow: 0 8px 48px rgba(0,0,0,0.18);
      flex-direction: column;
      overflow: hidden;
    `;

    // ── Modern Header ────────────────────────────────────────────────────────
    const header = document.createElement('div');
    header.style.cssText = `
      background: linear-gradient(135deg, ${settings.theme_color} 0%, ${settings.theme_color}cc 100%);
      color: white;
      padding: 13px 16px;
      display: flex;
      align-items: center;
      gap: 11px;
      flex-shrink: 0;
    `;
    const botName = settings.chatbot_name || 'Assistant';
    header.innerHTML = `
      <div style="width:38px;height:38px;border-radius:50%;background:rgba(255,255,255,0.2);display:flex;align-items:center;justify-content:center;flex-shrink:0;">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none"><path d="M12 2L13.09 9.26L20 12L13.09 14.74L12 22L10.91 14.74L4 12L10.91 9.26L12 2Z" fill="white"/><path d="M19 5L19.54 8.46L23 9L19.54 9.54L19 13L18.46 9.54L15 9L18.46 8.46L19 5Z" fill="rgba(255,255,255,0.7)"/></svg>
      </div>
      <div style="flex:1;min-width:0;">
        <div style="font-weight:700;font-size:14px;letter-spacing:0.2px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;">${botName}</div>
        <div style="font-size:11px;opacity:0.85;display:flex;align-items:center;gap:4px;margin-top:1px;">
          <span style="width:7px;height:7px;border-radius:50%;background:#4ade80;display:inline-block;box-shadow:0 0 4px #4ade80;"></span>Online
        </div>
      </div>
      <button id="chatbot-close-button" style="background:rgba(255,255,255,0.18);border:none;color:white;width:30px;height:30px;border-radius:50%;cursor:pointer;font-size:18px;display:flex;align-items:center;justify-content:center;transition:background 0.2s;flex-shrink:0;">&times;</button>
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
      padding: 16px 14px;
      background: #f0f4ff;
      scroll-behavior: smooth;
    `;

    // ── Modern input area ──
    const inputContainer = document.createElement('div');
    inputContainer.style.cssText = `
      display: flex;
      flex-direction: column;
      padding: 10px 14px 8px;
      background: white;
      box-shadow: 0 -2px 12px rgba(0,0,0,0.06);
      flex-shrink: 0;
    `;

    const inputRow = document.createElement('div');
    inputRow.style.cssText = 'display:flex;align-items:center;gap:10px;';

    const input = document.createElement('input');
    input.id = 'chatbot-input';
    input.type = 'text';
    input.placeholder = 'Ask me anything...';
    input.style.cssText = `
      flex: 1;
      padding: 11px 16px;
      border: 1.5px solid #e8eaf0;
      border-radius: 24px;
      outline: none;
      font-size: 14px;
      background: #f7f8fc;
      color: #1a1a2e;
      font-family: inherit;
      transition: border-color 0.2s, background 0.2s, box-shadow 0.2s;
    `;
    input.onfocus = () => {
      input.style.borderColor = settings.theme_color;
      input.style.background = '#fff';
      input.style.boxShadow = `0 0 0 3px ${settings.theme_color}22`;
    };
    input.onblur = () => {
      input.style.borderColor = '#e8eaf0';
      input.style.background = '#f7f8fc';
      input.style.boxShadow = 'none';
    };

    const sendButton = document.createElement('button');
    sendButton.id = 'chatbot-send-button';
    sendButton.innerHTML = `<svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M22 2L11 13" stroke="white" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"/><path d="M22 2L15 22L11 13L2 9L22 2Z" stroke="white" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"/></svg>`;
    sendButton.style.cssText = `
      background: linear-gradient(135deg, ${settings.theme_color}, ${settings.theme_color}cc);
      border: none;
      border-radius: 50%;
      width: 42px;
      height: 42px;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
      box-shadow: 0 3px 10px ${settings.theme_color}55;
      transition: transform 0.15s, box-shadow 0.15s;
    `;
    sendButton.onmouseover = () => { sendButton.style.transform = 'scale(1.08)'; sendButton.style.boxShadow = `0 5px 16px ${settings.theme_color}77`; };
    sendButton.onmouseout  = () => { sendButton.style.transform = 'scale(1)';    sendButton.style.boxShadow = `0 3px 10px ${settings.theme_color}55`; };
    sendButton.onclick = handleSendMessage;
    input.addEventListener('keypress', (e) => {
      if (e.key === 'Enter') handleSendMessage();
    });

    inputRow.appendChild(input);
    inputRow.appendChild(sendButton);

    // Disclaimer
    const disclaimer = document.createElement('div');
    disclaimer.style.cssText = 'text-align:center;padding-top:5px;font-size:11px;color:#b0b4c0;';
    const dBotName = settings.chatbot_name || 'AI';
    disclaimer.innerHTML = `<a href="#" onclick="return false;" style="color:${settings.theme_color};text-decoration:none;font-weight:600;">${dBotName}</a> can make mistakes. Verify important info.`;

    inputContainer.appendChild(inputRow);
    inputContainer.appendChild(disclaimer);

    wrapper.appendChild(messagesContainer);
    wrapper.appendChild(inputContainer);

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
      margin-bottom: 14px;
      display: flex;
      align-items: flex-end;
      gap: 8px;
      ${role === 'user' ? 'justify-content: flex-end;' : 'justify-content: flex-start;'}
    `;

    if (role === 'assistant') {
      // Bot avatar dot
      const avatarEl = document.createElement('div');
      if (settings.avatar_url) {
        const img = document.createElement('img');
        img.src = settings.avatar_url;
        img.style.cssText = 'width:30px;height:30px;border-radius:50%;object-fit:cover;flex-shrink:0;';
        img.onerror = () => {
          img.style.display = 'none';
          avatarEl.style.cssText += 'display:flex;';
        };
        avatarEl.appendChild(img);
        avatarEl.style.cssText = 'flex-shrink:0;';
      } else {
        avatarEl.style.cssText = `
          width:30px;height:30px;border-radius:50%;flex-shrink:0;
          background:linear-gradient(135deg,${settings.theme_color},${settings.theme_color}bb);
          display:flex;align-items:center;justify-content:center;
          box-shadow:0 2px 6px ${settings.theme_color}44;
        `;
        avatarEl.innerHTML = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M12 2L13.09 9.26L20 12L13.09 14.74L12 22L10.91 14.74L4 12L10.91 9.26L12 2Z" fill="white"/></svg>`;
      }
      messageDiv.appendChild(avatarEl);
    }

    const bubble = document.createElement('div');

    if (role === 'user') {
      bubble.style.cssText = `
        max-width: 72%;
        padding: 10px 15px;
        border-radius: 20px 20px 4px 20px;
        background: linear-gradient(135deg, ${settings.theme_color} 0%, ${settings.theme_color}dd 100%);
        color: white;
        font-size: 14px;
        font-weight: 500;
        line-height: 1.55;
        word-wrap: break-word;
        box-shadow: 0 3px 10px ${settings.theme_color}44;
      `;
      bubble.textContent = content;
    } else {
      bubble.style.cssText = `
        max-width: 80%;
        padding: 11px 14px;
        border-radius: 4px 20px 20px 20px;
        background: white;
        color: #1a1a2e;
        font-size: 14px;
        line-height: 1.6;
        word-wrap: break-word;
        box-shadow: 0 2px 10px rgba(0,0,0,0.08);
        border: 1px solid #eef0f8;
      `;
      bubble.classList.add('chatbot-response-bubble');
      bubble.innerHTML = formatMessageContent(content);
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
    typingDiv.style.cssText = 'margin-bottom:14px;display:flex;align-items:flex-end;gap:8px;justify-content:flex-start;';

    // Same avatar as messages
    const avatarEl = document.createElement('div');
    if (settings.avatar_url) {
      const img = document.createElement('img');
      img.src = settings.avatar_url;
      img.style.cssText = 'width:30px;height:30px;border-radius:50%;object-fit:cover;flex-shrink:0;';
      img.onerror = () => { img.style.display = 'none'; };
      avatarEl.appendChild(img);
      avatarEl.style.cssText = 'flex-shrink:0;';
    } else {
      avatarEl.style.cssText = `
        width:30px;height:30px;border-radius:50%;flex-shrink:0;
        background:linear-gradient(135deg,${settings.theme_color},${settings.theme_color}bb);
        display:flex;align-items:center;justify-content:center;
        box-shadow:0 2px 6px ${settings.theme_color}44;
      `;
      avatarEl.innerHTML = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M12 2L13.09 9.26L20 12L13.09 14.74L12 22L10.91 14.74L4 12L10.91 9.26L12 2Z" fill="white"/></svg>`;
    }
    typingDiv.appendChild(avatarEl);

    const bubble = document.createElement('div');
    bubble.style.cssText = `
      background: white;
      padding: 13px 16px;
      border-radius: 4px 20px 20px 20px;
      box-shadow: 0 2px 10px rgba(0,0,0,0.08);
      border: 1px solid #eef0f8;
      display: flex;
      align-items: center;
      gap: 5px;
    `;
    bubble.innerHTML = `
      <span style="width:8px;height:8px;border-radius:50%;background:${settings.theme_color};display:inline-block;animation:cb-bounce 1.2s infinite 0s;opacity:0.8;"></span>
      <span style="width:8px;height:8px;border-radius:50%;background:${settings.theme_color};display:inline-block;animation:cb-bounce 1.2s infinite 0.2s;opacity:0.8;"></span>
      <span style="width:8px;height:8px;border-radius:50%;background:${settings.theme_color};display:inline-block;animation:cb-bounce 1.2s infinite 0.4s;opacity:0.8;"></span>
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

  // ── Load Google Sans from Google Fonts (only once) ─────────────────────────
  if (!document.getElementById('chatbot-google-sans-font')) {
    const preconnect1 = document.createElement('link');
    preconnect1.rel = 'preconnect';
    preconnect1.href = 'https://fonts.googleapis.com';
    document.head.appendChild(preconnect1);

    const preconnect2 = document.createElement('link');
    preconnect2.rel = 'preconnect';
    preconnect2.href = 'https://fonts.gstatic.com';
    preconnect2.crossOrigin = 'anonymous';
    document.head.appendChild(preconnect2);

    const fontLink = document.createElement('link');
    fontLink.id = 'chatbot-google-sans-font';
    fontLink.rel = 'stylesheet';
    fontLink.href = 'https://fonts.googleapis.com/css2?family=Google+Sans:ital,opsz,wght@0,17..18,400..700;1,17..18,400..700&display=swap';
    document.head.appendChild(fontLink);
  }

  // Styles (theme-specific keyframes are appended in createWidget once settings load)
  const style = document.createElement('style');
  style.textContent = `
    @keyframes blink {
      0%, 50%, 100% { opacity: 0.3; }
      25% { opacity: 1; }
    }
    @keyframes cb-bounce {
      0%, 60%, 100% { transform: translateY(0); }
      30% { transform: translateY(-7px); }
    }
    #chatbot-widget-container,
    #chatbot-widget-container * {
      font-family: 'Google Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif !important;
      font-optical-sizing: auto;
      -webkit-font-smoothing: antialiased;
      -moz-osx-font-smoothing: grayscale;
      box-sizing: border-box;
    }
    #chatbot-messages::-webkit-scrollbar { width: 4px; }
    #chatbot-messages::-webkit-scrollbar-track { background: transparent; }
    #chatbot-messages::-webkit-scrollbar-thumb { background: #d0d4e8; border-radius: 4px; }
    #chatbot-close-button:hover { background: rgba(255,255,255,0.3) !important; }
  `;
  document.head.appendChild(style);

  // Boot
  fetchSettings().then(() => {
    createWidget();
  });
})();
