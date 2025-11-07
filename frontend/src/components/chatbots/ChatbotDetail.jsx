import React, { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Paper,
  Typography,
  TextField,
  Button,
  Grid,
  Alert,
  Divider,
  Card,
  CardContent,
  Switch,
  FormControlLabel,
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  Visibility as VisibilityIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material';
import api from '../../services/api';

function ChatbotDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [chatbot, setChatbot] = useState(null);
  const [settings, setSettings] = useState(null);
  const [loading, setLoading] = useState(true);
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');
  const [showPreview, setShowPreview] = useState(true);
  const [previewKey, setPreviewKey] = useState(0);
  const previewRef = useRef(null);

  useEffect(() => {
    fetchChatbotData();
  }, [id]);

  const fetchChatbotData = async () => {
    try {
      const [chatbotRes, settingsRes] = await Promise.all([
        api.get(`/chatbots/${id}`),
        api.get(`/chatbots/${id}/settings`),
      ]);
      setChatbot(chatbotRes.data);
      setSettings(settingsRes.data);
    } catch (error) {
      console.error('Failed to fetch chatbot data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateSettings = async () => {
    setError('');
    setSuccess('');

    try {
      await api.put(`/chatbots/${id}/settings`, settings);
      setSuccess('Settings updated successfully');
      // Refresh preview after save
      setTimeout(() => refreshPreview(), 500);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update settings');
    }
  };

  const refreshPreview = () => {
    setPreviewKey(prev => prev + 1);
  };

  const handleSettingChange = (field, value) => {
    setSettings({ ...settings, [field]: value });
    // Auto-refresh preview after a short delay
    setTimeout(() => refreshPreview(), 300);
  };

  const widgetCode = `<script>
  (function() {
    var script = document.createElement('script');
    script.src = 'http://localhost:8081/widget.js';
    script.setAttribute('data-chatbot-id', '${id}');
    document.body.appendChild(script);
  })();
</script>`;

  if (loading) {
    return <Typography>Loading...</Typography>;
  }

  return (
    <Box>
      <Button startIcon={<ArrowBackIcon />} onClick={() => navigate('/chatbots')} sx={{ mb: 2 }}>
        Back to Chatbots
      </Button>

      <Typography variant="h4" gutterBottom>
        {chatbot?.name}
      </Typography>

      {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}
      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6">
                Chatbot Settings
              </Typography>
              <FormControlLabel
                control={
                  <Switch
                    checked={showPreview}
                    onChange={(e) => setShowPreview(e.target.checked)}
                    color="primary"
                  />
                }
                label="Live Preview"
              />
            </Box>
            <TextField
              fullWidth
              label="Theme Color"
              value={settings?.theme_color || '#000000'}
              onChange={(e) => handleSettingChange('theme_color', e.target.value)}
              margin="normal"
              type="color"
              helperText="Choose your brand color"
            />
            <TextField
              fullWidth
              label="Position"
              value={settings?.position || 'bottom-left'}
              onChange={(e) => handleSettingChange('position', e.target.value)}
              margin="normal"
              select
              SelectProps={{ native: true }}
            >
              <option value="bottom-right">Bottom Right</option>
              <option value="bottom-left">Bottom Left</option>
              <option value="top-right">Top Right</option>
              <option value="top-left">Top Left</option>
            </TextField>
            <TextField
              fullWidth
              label="Welcome Message"
              value={settings?.welcome_message || ''}
              onChange={(e) => handleSettingChange('welcome_message', e.target.value)}
              margin="normal"
              multiline
              rows={2}
              helperText="First message users see"
            />
            <TextField
              fullWidth
              label="Avatar URL"
              value={settings?.avatar_url || ''}
              onChange={(e) => handleSettingChange('avatar_url', e.target.value)}
              margin="normal"
              helperText="Optional bot avatar image"
            />
            <TextField
              fullWidth
              label="Widget Size"
              value={settings?.widget_size || 'medium'}
              onChange={(e) => handleSettingChange('widget_size', e.target.value)}
              margin="normal"
              select
              SelectProps={{ native: true }}
            >
              <option value="small">Small</option>
              <option value="medium">Medium</option>
              <option value="large">Large</option>
            </TextField>
            <Button
              variant="contained"
              onClick={handleUpdateSettings}
              sx={{ mt: 2 }}
              fullWidth
            >
              Save Settings
            </Button>
          </Paper>
        </Grid>

        <Grid item xs={12} md={6}>
          {/* Live Preview */}
          {showPreview && (
            <Paper sx={{ p: 3, mb: 3, minHeight: 600 }}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <VisibilityIcon /> Live Preview
                </Typography>
                <Button
                  size="small"
                  startIcon={<RefreshIcon />}
                  onClick={refreshPreview}
                  variant="outlined"
                >
                  Refresh
                </Button>
              </Box>
              <Alert severity="info" sx={{ mb: 2 }}>
                Interactive preview - Click the widget icon to test conversations!
              </Alert>
              <Box
                sx={{
                  position: 'relative',
                  width: '100%',
                  height: 700,
                  border: '2px solid #e0e0e0',
                  borderRadius: 2,
                  overflow: 'hidden',
                }}
              >
                <iframe
                  key={previewKey}
                  ref={previewRef}
                  srcDoc={`
                    <!DOCTYPE html>
                    <html>
                    <head>
                      <meta charset="UTF-8">
                      <meta name="viewport" content="width=device-width, initial-scale=1.0">
                      <style>
                        body {
                          margin: 0;
                          padding: 0;
                          font-family: Arial, sans-serif;
                          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
                          min-height: 100vh;
                          display: flex;
                          align-items: center;
                          justify-content: center;
                          color: white;
                        }
                        .info {
                          text-align: center;
                          padding: 40px;
                        }
                        .info h2 {
                          margin-bottom: 20px;
                          font-size: 24px;
                        }
                        .info p {
                          margin: 10px 0;
                          opacity: 0.9;
                        }
                        .badge {
                          display: inline-block;
                          background: rgba(255,255,255,0.2);
                          padding: 5px 15px;
                          border-radius: 20px;
                          margin: 5px;
                          font-size: 14px;
                        }
                      </style>
                    </head>
                    <body>
                      <div class="info">
                        <h2>🤖 Live Chatbot Preview</h2>
                        <p>Your chatbot widget will appear in the <strong>${settings?.position || 'bottom-left'}</strong> corner</p>
                        <div>
                          <span class="badge">Position: ${settings?.position || 'bottom-left'}</span>
                          <span class="badge">Size: ${settings?.widget_size || 'medium'}</span>
                        </div>
                        <p style="margin-top: 30px; font-size: 14px;">
                          👇 Look for the chat icon below and click it to test!
                        </p>
                      </div>
                      <script>
                        (function() {
                          var script = document.createElement('script');
                          script.src = 'http://localhost:8081/widget.js';
                          script.setAttribute('data-chatbot-id', '${id}');
                          document.body.appendChild(script);
                        })();
                      </script>
                    </body>
                    </html>
                  `}
                  style={{
                    width: '100%',
                    height: '100%',
                    border: 'none',
                  }}
                  title="Chatbot Preview"
                />
              </Box>
            </Paper>
          )}
        </Grid>

        <Grid item xs={12} md={showPreview ? 6 : 12}>
          <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom>
              Widget Code
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              Copy and paste this code into your website's HTML:
            </Typography>
            <TextField
              fullWidth
              multiline
              rows={8}
              value={widgetCode}
              margin="normal"
              InputProps={{ readOnly: true }}
            />
            <Button
              variant="outlined"
              onClick={() => navigator.clipboard.writeText(widgetCode)}
              sx={{ mt: 1 }}
            >
              Copy Code
            </Button>
          </Paper>

          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Quick Actions
              </Typography>
              <Button
                variant="outlined"
                fullWidth
                sx={{ mb: 1 }}
                onClick={() => navigate(`/knowledge/${id}`)}
              >
                Manage Knowledge Base
              </Button>
              <Button
                variant="outlined"
                fullWidth
                onClick={() => navigate(`/analytics/${id}`)}
              >
                View Analytics
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
}

export default ChatbotDetail;

