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
  Card,
  CardContent,
  Switch,
  FormControlLabel,
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  IconButton,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Divider,
  Dialog,
  DialogTitle,
  DialogContent,
  CircularProgress,
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  Visibility as VisibilityIcon,
  Refresh as RefreshIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
  People as PeopleIcon,
} from '@mui/icons-material';
import api from '../../services/api';

const DEFAULT_FIELDS = [
  { name: 'name', label: 'Your Name', type: 'text', required: true, placeholder: 'Enter your name' },
  { name: 'email', label: 'Email Address', type: 'email', required: true, placeholder: 'Enter your email' },
];

function TabPanel({ children, value, index }) {
  return value === index ? <Box sx={{ pt: 3 }}>{children}</Box> : null;
}

function ChatbotDetail() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [activeTab, setActiveTab] = useState(0);
  const [chatbot, setChatbot] = useState(null);
  const [settings, setSettings] = useState(null);
  const [loading, setLoading] = useState(true);
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');
  const [showPreview, setShowPreview] = useState(true);
  const [previewKey, setPreviewKey] = useState(0);
  const previewRef = useRef(null);

  // Lead capture state
  const [leadConfig, setLeadConfig] = useState({
    enabled: false,
    title: 'Before we begin...',
    subtitle: 'Please share a few details so we can assist you better.',
    fields: DEFAULT_FIELDS,
  });
  const [leadSuccess, setLeadSuccess] = useState('');
  const [leadError, setLeadError] = useState('');

  // Leads list state
  const [leads, setLeads] = useState([]);
  const [leadsLoading, setLeadsLoading] = useState(false);

  // Session chat dialog
  const [sessionDialog, setSessionDialog] = useState({ open: false, lead: null, messages: [], loading: false });

  // Suggestion chips state
  const [newSuggestion, setNewSuggestion] = useState('');

  useEffect(() => {
    fetchChatbotData();
  }, [id]);

  useEffect(() => {
    if (activeTab === 2) fetchLeads();
  }, [activeTab]);

  const fetchChatbotData = async () => {
    try {
      const [chatbotRes, settingsRes, leadRes] = await Promise.all([
        api.get(`/chatbots/${id}`),
        api.get(`/chatbots/${id}/settings`),
        api.get(`/chatbots/${id}/lead-capture`),
      ]);
      setChatbot(chatbotRes.data);
      setSettings({
        ...settingsRes.data,
        suggestions: Array.isArray(settingsRes.data.suggestions) ? settingsRes.data.suggestions : [],
      });
      setLeadConfig({
        enabled: leadRes.data.enabled || false,
        title: leadRes.data.title || 'Before we begin...',
        subtitle: leadRes.data.subtitle || '',
        fields: leadRes.data.fields?.length ? leadRes.data.fields : DEFAULT_FIELDS,
      });
    } catch (err) {
      console.error('Failed to fetch chatbot data:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchLeads = async () => {
    setLeadsLoading(true);
    try {
      const res = await api.get(`/chatbots/${id}/leads`);
      setLeads(res.data || []);
    } catch (err) {
      console.error('Failed to fetch leads:', err);
    } finally {
      setLeadsLoading(false);
    }
  };

  const openSession = async (lead) => {
    setSessionDialog({ open: true, lead, messages: [], loading: true });
    try {
      const res = await api.get(`/chatbots/${id}/sessions/${lead.session_id}/messages`);
      setSessionDialog(prev => ({ ...prev, messages: res.data.messages || [], loading: false }));
    } catch (err) {
      setSessionDialog(prev => ({ ...prev, loading: false }));
    }
  };

  const closeSession = () => setSessionDialog({ open: false, lead: null, messages: [], loading: false });

  // ── Settings tab ────────────────────────────────────────────────────────────

  const handleUpdateSettings = async () => {
    setError(''); setSuccess('');
    try {
      await api.put(`/chatbots/${id}/settings`, settings);
      setSuccess('Settings updated successfully');
      setTimeout(() => refreshPreview(), 500);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update settings');
    }
  };

  const refreshPreview = () => setPreviewKey(prev => prev + 1);

  const handleSettingChange = (field, value) => {
    setSettings({ ...settings, [field]: value });
    setTimeout(() => refreshPreview(), 300);
  };

  // ── Lead Capture tab ────────────────────────────────────────────────────────

  const handleSaveLeadConfig = async () => {
    setLeadError(''); setLeadSuccess('');
    try {
      await api.put(`/chatbots/${id}/lead-capture`, leadConfig);
      setLeadSuccess('Lead capture settings saved successfully');
    } catch (err) {
      setLeadError(err.response?.data?.error || 'Failed to save lead capture settings');
    }
  };

  const addField = () => {
    setLeadConfig(prev => ({
      ...prev,
      fields: [...prev.fields, { name: '', label: '', type: 'text', required: false, placeholder: '' }],
    }));
  };

  const removeField = (index) => {
    setLeadConfig(prev => ({
      ...prev,
      fields: prev.fields.filter((_, i) => i !== index),
    }));
  };

  const updateField = (index, key, value) => {
    setLeadConfig(prev => {
      const fields = [...prev.fields];
      fields[index] = { ...fields[index], [key]: value };
      return { ...prev, fields };
    });
  };

  // ── Widget helpers ──────────────────────────────────────────────────────────

  const getWidgetUrl = () => {
    const envUrl = import.meta.env.VITE_WIDGET_URL;
    if (envUrl) return envUrl;
    const { hostname, protocol } = window.location;
    if (hostname === 'localhost' || hostname === '127.0.0.1') return 'http://localhost:8081/widget.js';
    return `${protocol}//chatbot-api.appster.co.in/widget.js`;
  };

  const getApiUrl = () => {
    const { hostname, protocol } = window.location;
    if (hostname === 'localhost' || hostname === '127.0.0.1') return 'http://localhost:8081/api';
    return `${protocol}//chatbot-api.appster.co.in/api`;
  };

  const widgetCode = `<script>
  (function() {
    var script = document.createElement('script');
    script.src = '${getWidgetUrl()}';
    script.setAttribute('data-chatbot-id', '${id}');
    document.body.appendChild(script);
  })();
</script>`;

  if (loading) return <Typography>Loading...</Typography>;

  // ── Leads table helpers ─────────────────────────────────────────────────────

  // Collect all unique field keys across all leads for table columns
  const leadFieldKeys = Array.from(
    new Set(leads.flatMap(l => Object.keys(l.field_values || {})))
  );

  return (
    <Box>
      <Button startIcon={<ArrowBackIcon />} onClick={() => navigate('/chatbots')} sx={{ mb: 2 }}>
        Back to Chatbots
      </Button>

      <Typography variant="h4" gutterBottom>{chatbot?.name}</Typography>

      <Paper sx={{ mb: 3 }}>
        <Tabs value={activeTab} onChange={(_, v) => setActiveTab(v)} sx={{ borderBottom: 1, borderColor: 'divider', px: 2 }}>
          <Tab label="Settings" />
          <Tab label="Lead Capture" />
          <Tab label={
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
              <PeopleIcon fontSize="small" /> Leads {leads.length > 0 && <Chip label={leads.length} size="small" color="primary" />}
            </Box>
          } />
        </Tabs>

        <Box sx={{ p: 3 }}>

          {/* ── Tab 0: Settings ───────────────────────────────────────────── */}
          <TabPanel value={activeTab} index={0}>
            {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}
            {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6">Chatbot Settings</Typography>
                  <FormControlLabel
                    control={<Switch checked={showPreview} onChange={e => setShowPreview(e.target.checked)} color="primary" />}
                    label="Live Preview"
                  />
                </Box>
                <TextField fullWidth label="Theme Color" value={settings?.theme_color || '#000000'} onChange={e => handleSettingChange('theme_color', e.target.value)} margin="normal" type="color" helperText="Choose your brand color" />
                <TextField fullWidth label="Position" value={settings?.position || 'bottom-right'} onChange={e => handleSettingChange('position', e.target.value)} margin="normal" select SelectProps={{ native: true }}>
                  <option value="bottom-right">Bottom Right</option>
                  <option value="bottom-left">Bottom Left</option>
                </TextField>
                <TextField fullWidth label="Welcome Message" value={settings?.welcome_message || ''} onChange={e => handleSettingChange('welcome_message', e.target.value)} margin="normal" multiline rows={2} helperText="First message users see" />
                <TextField fullWidth label="Avatar URL" value={settings?.avatar_url || ''} onChange={e => handleSettingChange('avatar_url', e.target.value)} margin="normal" helperText="Optional bot avatar image" />
                <TextField fullWidth label="Widget Size" value={settings?.widget_size || 'medium'} onChange={e => handleSettingChange('widget_size', e.target.value)} margin="normal" select SelectProps={{ native: true }}>
                  <option value="small">Small</option>
                  <option value="medium">Medium</option>
                  <option value="large">Large</option>
                </TextField>
                {/* Quick Suggestions */}
                <Divider sx={{ my: 3 }} />
                <Typography variant="h6" gutterBottom>Quick Suggestions</Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  Clickable chips shown in the chat after the welcome message. Users can click them to instantly send that message.
                </Typography>
                <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Add a suggestion"
                    placeholder="e.g. What are your fees?"
                    value={newSuggestion}
                    onChange={e => setNewSuggestion(e.target.value)}
                    onKeyDown={e => {
                      if (e.key === 'Enter' && newSuggestion.trim()) {
                        handleSettingChange('suggestions', [...(settings?.suggestions || []), newSuggestion.trim()]);
                        setNewSuggestion('');
                      }
                    }}
                  />
                  <Button
                    variant="outlined"
                    startIcon={<AddIcon />}
                    onClick={() => {
                      if (!newSuggestion.trim()) return;
                      handleSettingChange('suggestions', [...(settings?.suggestions || []), newSuggestion.trim()]);
                      setNewSuggestion('');
                    }}
                    sx={{ whiteSpace: 'nowrap' }}
                  >
                    Add
                  </Button>
                </Box>
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, minHeight: 36 }}>
                  {(settings?.suggestions || []).map((s, i) => (
                    <Chip
                      key={i}
                      label={s}
                      onDelete={() => handleSettingChange('suggestions', settings.suggestions.filter((_, idx) => idx !== i))}
                      color="primary"
                      variant="outlined"
                      size="small"
                    />
                  ))}
                  {(settings?.suggestions || []).length === 0 && (
                    <Typography variant="body2" color="text.disabled" sx={{ fontStyle: 'italic' }}>No suggestions added yet</Typography>
                  )}
                </Box>

                <Button variant="contained" onClick={handleUpdateSettings} sx={{ mt: 3 }} fullWidth>Save Settings</Button>

                <Divider sx={{ my: 3 }} />

                <Typography variant="h6" gutterBottom>Widget Code</Typography>
                <Typography variant="body2" color="text.secondary" gutterBottom>Copy and paste into your website's HTML:</Typography>
                <TextField fullWidth multiline rows={6} value={widgetCode} margin="normal" InputProps={{ readOnly: true }} />
                <Button variant="outlined" onClick={() => navigator.clipboard.writeText(widgetCode)} sx={{ mt: 1 }}>Copy Code</Button>

                <Divider sx={{ my: 3 }} />

                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>Quick Actions</Typography>
                    <Button variant="outlined" fullWidth sx={{ mb: 1 }} onClick={() => navigate(`/knowledge/${id}`)}>Manage Knowledge Base</Button>
                    <Button variant="outlined" fullWidth onClick={() => navigate(`/analytics/${id}`)}>View Analytics</Button>
                  </CardContent>
                </Card>
              </Grid>

              <Grid item xs={12} md={6}>
                {showPreview && (
                  <Box>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                      <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <VisibilityIcon /> Live Preview
                      </Typography>
                      <Button size="small" startIcon={<RefreshIcon />} onClick={refreshPreview} variant="outlined">Refresh</Button>
                    </Box>
                    <Alert severity="info" sx={{ mb: 2 }}>Interactive preview — click the widget icon to test!</Alert>
                    <Box sx={{ position: 'relative', width: '100%', height: 700, border: '2px solid #e0e0e0', borderRadius: 2, overflow: 'hidden' }}>
                      <iframe
                        key={previewKey}
                        ref={previewRef}
                        srcDoc={`<!DOCTYPE html><html><head><meta charset="UTF-8"><style>body{margin:0;padding:0;font-family:Arial,sans-serif;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);min-height:100vh;display:flex;align-items:center;justify-content:center;color:white;}.info{text-align:center;padding:40px;}</style></head><body><div class="info"><h2>🤖 Live Chatbot Preview</h2><p>Your chatbot widget will appear in the <strong>${settings?.position || 'bottom-right'}</strong> corner</p><p style="margin-top:20px;font-size:14px;">👇 Look for the chat icon and click it to test!</p></div><script>(function(){var s=document.createElement('script');s.src='${getWidgetUrl()}';s.setAttribute('data-chatbot-id','${id}');s.setAttribute('data-api-url','${getApiUrl()}');document.body.appendChild(s);})();</script></body></html>`}
                        style={{ width: '100%', height: '100%', border: 'none' }}
                        title="Chatbot Preview"
                      />
                    </Box>
                  </Box>
                )}
              </Grid>
            </Grid>
          </TabPanel>

          {/* ── Tab 1: Lead Capture ───────────────────────────────────────── */}
          <TabPanel value={activeTab} index={1}>
            {leadSuccess && <Alert severity="success" sx={{ mb: 2 }}>{leadSuccess}</Alert>}
            {leadError && <Alert severity="error" sx={{ mb: 2 }}>{leadError}</Alert>}

            <Grid container spacing={3}>
              <Grid item xs={12} md={7}>
                <Paper variant="outlined" sx={{ p: 3 }}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={leadConfig.enabled}
                        onChange={e => setLeadConfig(prev => ({ ...prev, enabled: e.target.checked }))}
                        color="primary"
                      />
                    }
                    label={<Typography fontWeight={600}>Enable Lead Capture Form</Typography>}
                  />
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 3, mt: 0.5 }}>
                    When enabled, visitors must fill this form before they can start chatting.
                  </Typography>

                  <TextField
                    fullWidth
                    label="Form Title"
                    value={leadConfig.title}
                    onChange={e => setLeadConfig(prev => ({ ...prev, title: e.target.value }))}
                    margin="normal"
                    helperText='e.g. "Before we begin..."'
                  />
                  <TextField
                    fullWidth
                    label="Form Subtitle"
                    value={leadConfig.subtitle}
                    onChange={e => setLeadConfig(prev => ({ ...prev, subtitle: e.target.value }))}
                    margin="normal"
                    multiline
                    rows={2}
                    helperText="Optional description shown below the title"
                  />

                  <Divider sx={{ my: 3 }} />

                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                    <Typography variant="subtitle1" fontWeight={600}>Form Fields</Typography>
                    <Button startIcon={<AddIcon />} size="small" variant="outlined" onClick={addField}>Add Field</Button>
                  </Box>

                  {leadConfig.fields.map((field, idx) => (
                    <Paper key={idx} variant="outlined" sx={{ p: 2, mb: 2, bgcolor: '#fafafa' }}>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                        <Typography variant="body2" color="text.secondary" fontWeight={500}>Field {idx + 1}</Typography>
                        <IconButton size="small" color="error" onClick={() => removeField(idx)} disabled={leadConfig.fields.length <= 1}>
                          <DeleteIcon fontSize="small" />
                        </IconButton>
                      </Box>
                      <Grid container spacing={1.5}>
                        <Grid item xs={6}>
                          <TextField
                            fullWidth size="small" label="Field Name (key)"
                            value={field.name}
                            onChange={e => updateField(idx, 'name', e.target.value)}
                            helperText='e.g. "email"'
                          />
                        </Grid>
                        <Grid item xs={6}>
                          <TextField
                            fullWidth size="small" label="Display Label"
                            value={field.label}
                            onChange={e => updateField(idx, 'label', e.target.value)}
                            helperText='e.g. "Email Address"'
                          />
                        </Grid>
                        <Grid item xs={6}>
                          <FormControl fullWidth size="small">
                            <InputLabel>Type</InputLabel>
                            <Select label="Type" value={field.type} onChange={e => updateField(idx, 'type', e.target.value)}>
                              <MenuItem value="text">Text</MenuItem>
                              <MenuItem value="email">Email</MenuItem>
                              <MenuItem value="tel">Phone</MenuItem>
                              <MenuItem value="textarea">Textarea</MenuItem>
                            </Select>
                          </FormControl>
                        </Grid>
                        <Grid item xs={6}>
                          <TextField
                            fullWidth size="small" label="Placeholder"
                            value={field.placeholder}
                            onChange={e => updateField(idx, 'placeholder', e.target.value)}
                          />
                        </Grid>
                        <Grid item xs={12}>
                          <FormControlLabel
                            control={<Switch size="small" checked={field.required} onChange={e => updateField(idx, 'required', e.target.checked)} />}
                            label={<Typography variant="body2">Required</Typography>}
                          />
                        </Grid>
                      </Grid>
                    </Paper>
                  ))}

                  <Button variant="contained" onClick={handleSaveLeadConfig} fullWidth sx={{ mt: 1 }}>
                    Save Lead Capture Settings
                  </Button>
                </Paper>
              </Grid>

              <Grid item xs={12} md={5}>
                <Paper variant="outlined" sx={{ p: 3, bgcolor: '#f5f5f5' }}>
                  <Typography variant="subtitle1" fontWeight={600} gutterBottom>Preview</Typography>
                  <Divider sx={{ mb: 2 }} />
                  <Typography fontWeight={600} fontSize={15} sx={{ mb: 0.5 }}>{leadConfig.title}</Typography>
                  {leadConfig.subtitle && (
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>{leadConfig.subtitle}</Typography>
                  )}
                  {leadConfig.fields.map((f, i) => (
                    <Box key={i} sx={{ mb: 1.5 }}>
                      <Typography variant="caption" fontWeight={500} color="text.secondary">
                        {f.label}{f.required ? ' *' : ''}
                      </Typography>
                      <Box sx={{ mt: 0.5, p: 1, border: '1px solid #ddd', borderRadius: 1, bgcolor: 'white', fontSize: 13, color: '#999' }}>
                        {f.placeholder || f.label}
                      </Box>
                    </Box>
                  ))}
                  <Box sx={{ mt: 2, p: 1.5, bgcolor: settings?.theme_color || '#1976d2', color: 'white', borderRadius: 1, textAlign: 'center', fontSize: 14, fontWeight: 600 }}>
                    Start Chat
                  </Box>
                </Paper>
              </Grid>
            </Grid>
          </TabPanel>

          {/* ── Tab 2: Leads ──────────────────────────────────────────────── */}
          <TabPanel value={activeTab} index={2}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6">Captured Leads</Typography>
              <Button startIcon={<RefreshIcon />} size="small" variant="outlined" onClick={fetchLeads}>Refresh</Button>
            </Box>

            {leadsLoading ? (
              <Typography color="text.secondary">Loading leads...</Typography>
            ) : leads.length === 0 ? (
              <Alert severity="info">No leads captured yet. Enable the lead capture form and share your chatbot widget.</Alert>
            ) : (
              <TableContainer component={Paper} variant="outlined">
                <Table size="small">
                  <TableHead>
                    <TableRow sx={{ bgcolor: '#f5f5f5' }}>
                      <TableCell><strong>#</strong></TableCell>
                      <TableCell><strong>Date</strong></TableCell>
                      {leadFieldKeys.map(key => (
                        <TableCell key={key}><strong>{key.charAt(0).toUpperCase() + key.slice(1)}</strong></TableCell>
                      ))}
                      <TableCell><strong>Session</strong></TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {leads.map((lead, idx) => (
                      <TableRow
                        key={lead.id}
                        hover
                        onClick={() => openSession(lead)}
                        sx={{ cursor: 'pointer' }}
                      >
                        <TableCell>{leads.length - idx}</TableCell>
                        <TableCell sx={{ whiteSpace: 'nowrap' }}>
                          {new Date(lead.created_at).toLocaleString()}
                        </TableCell>
                        {leadFieldKeys.map(key => (
                          <TableCell key={key}>{lead.field_values?.[key] || '—'}</TableCell>
                        ))}
                        <TableCell>
                          <Typography variant="caption" color="text.secondary" sx={{ fontFamily: 'monospace' }}>
                            {lead.session_id?.slice(0, 16)}…
                          </Typography>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            )}
          </TabPanel>

        </Box>
      </Paper>
      {/* ── Session Chat Dialog ──────────────────────────────────────────── */}
      <Dialog open={sessionDialog.open} onClose={closeSession} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ pb: 1 }}>
          <Typography fontWeight={700} fontSize={16}>Session Chat</Typography>
          {sessionDialog.lead && (
            <Box sx={{ mt: 0.5, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
              {Object.entries(sessionDialog.lead.field_values || {}).map(([k, v]) => (
                <Chip key={k} label={`${k}: ${v}`} size="small" variant="outlined" />
              ))}
              <Chip
                label={new Date(sessionDialog.lead.created_at).toLocaleString()}
                size="small"
                sx={{ bgcolor: '#f5f5f5' }}
              />
            </Box>
          )}
        </DialogTitle>

        <DialogContent dividers sx={{ p: 0 }}>
          {sessionDialog.loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
              <CircularProgress size={32} />
            </Box>
          ) : sessionDialog.messages.length === 0 ? (
            <Box sx={{ p: 3 }}>
              <Alert severity="info">No messages found for this session.</Alert>
            </Box>
          ) : (
            <Box sx={{ p: 2, display: 'flex', flexDirection: 'column', gap: 1.5, maxHeight: 480, overflowY: 'auto', bgcolor: '#f5f5f5' }}>
              {sessionDialog.messages.map(msg => (
                <Box
                  key={msg.id}
                  sx={{
                    display: 'flex',
                    justifyContent: msg.role === 'user' ? 'flex-end' : 'flex-start',
                  }}
                >
                  <Box
                    sx={{
                      maxWidth: '75%',
                      px: 2,
                      py: 1,
                      borderRadius: 2.5,
                      fontSize: 14,
                      lineHeight: 1.5,
                      wordBreak: 'break-word',
                      ...(msg.role === 'user'
                        ? { bgcolor: settings?.theme_color || '#1976d2', color: 'white' }
                        : { bgcolor: 'white', border: '1px solid #e0e0e0', color: '#333' }),
                    }}
                  >
                    {msg.content}
                    <Typography
                      variant="caption"
                      sx={{
                        display: 'block',
                        mt: 0.5,
                        opacity: 0.65,
                        fontSize: 11,
                        textAlign: msg.role === 'user' ? 'right' : 'left',
                      }}
                    >
                      {new Date(msg.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                    </Typography>
                  </Box>
                </Box>
              ))}
            </Box>
          )}
        </DialogContent>
      </Dialog>
    </Box>
  );
}

export default ChatbotDetail;
