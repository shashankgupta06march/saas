import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Paper,
  Typography,
  TextField,
  Button,
  List,
  ListItem,
  ListItemText,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  Tabs,
  Tab,
  Chip,
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  Delete as DeleteIcon,
  Add as AddIcon,
  CloudUpload as CloudUploadIcon,
  Link as LinkIcon,
  TextFields as TextFieldsIcon,
} from '@mui/icons-material';
import api from '../../services/api';

function KnowledgeBase() {
  const { chatbotId } = useParams();
  const navigate = useNavigate();
  const [knowledge, setKnowledge] = useState([]);
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [url, setUrl] = useState('');
  const [file, setFile] = useState(null);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);
  const [tabValue, setTabValue] = useState(0);

  useEffect(() => {
    fetchKnowledge();
  }, [chatbotId]);

  const fetchKnowledge = async () => {
    try {
      const response = await api.get(`/knowledge/chatbot/${chatbotId}`);
      setKnowledge(response.data || []);
    } catch (error) {
      console.error('Failed to fetch knowledge:', error);
    }
  };

  const handleAddText = async () => {
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      await api.post('/knowledge', {
        chatbot_id: parseInt(chatbotId),
        title,
        content,
        content_type: 'text',
      });
      setSuccess('Knowledge added successfully');
      closeDialog();
      fetchKnowledge();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to add knowledge');
    } finally {
      setLoading(false);
    }
  };

  const handleUploadFile = async () => {
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('chatbot_id', chatbotId);
      formData.append('title', title || file.name);

      await api.post('/knowledge/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      setSuccess('File uploaded and processed successfully');
      closeDialog();
      fetchKnowledge();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to upload file');
    } finally {
      setLoading(false);
    }
  };

  const handleScrapeURL = async () => {
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      await api.post('/knowledge/scrape', {
        chatbot_id: parseInt(chatbotId),
        url,
        title: title || url,
      });
      setSuccess('Website scraped and processed successfully');
      closeDialog();
      fetchKnowledge();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to scrape website');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this knowledge entry?')) {
      return;
    }

    try {
      await api.delete(`/knowledge/${id}`);
      setSuccess('Knowledge deleted successfully');
      fetchKnowledge();
    } catch (error) {
      setError('Failed to delete knowledge');
    }
  };

  const closeDialog = () => {
    setOpen(false);
    setTitle('');
    setContent('');
    setUrl('');
    setFile(null);
    setTabValue(0);
  };

  const handleFileChange = (event) => {
    const selectedFile = event.target.files[0];
    if (selectedFile) {
      const allowedTypes = ['.pdf', '.docx', '.txt'];
      const fileExt = '.' + selectedFile.name.split('.').pop().toLowerCase();
      
      if (!allowedTypes.includes(fileExt)) {
        setError('Only PDF, DOCX, and TXT files are supported');
        return;
      }
      
      setFile(selectedFile);
      if (!title) {
        setTitle(selectedFile.name);
      }
    }
  };

  const getContentTypeIcon = (contentType) => {
    switch (contentType) {
      case 'pdf':
        return <Chip label="PDF" size="small" color="error" sx={{ mr: 1 }} />;
      case 'docx':
        return <Chip label="DOCX" size="small" color="primary" sx={{ mr: 1 }} />;
      case 'webpage':
        return <Chip label="Webpage" size="small" color="info" sx={{ mr: 1 }} />;
      case 'text':
      default:
        return <Chip label="Text" size="small" color="success" sx={{ mr: 1 }} />;
    }
  };

  const canSubmit = () => {
    if (tabValue === 0) return title && content;
    if (tabValue === 1) return file;
    if (tabValue === 2) return url;
    return false;
  };

  const handleSubmit = () => {
    if (tabValue === 0) handleAddText();
    else if (tabValue === 1) handleUploadFile();
    else if (tabValue === 2) handleScrapeURL();
  };

  return (
    <Box>
      <Button startIcon={<ArrowBackIcon />} onClick={() => navigate('/chatbots')} sx={{ mb: 2 }}>
        Back to Chatbots
      </Button>

      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
        <Typography variant="h4">Knowledge Base</Typography>
        <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpen(true)}>
          Add Knowledge
        </Button>
      </Box>

      {success && <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccess('')}>{success}</Alert>}
      {error && <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>{error}</Alert>}

      <Paper>
        {knowledge.length === 0 ? (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography color="text.secondary">
              No knowledge entries yet. Add your first entry to train the chatbot!
            </Typography>
          </Box>
        ) : (
          <List>
            {knowledge.map((item, index) => (
              <ListItem
                key={item.id}
                secondaryAction={
                  <IconButton edge="end" onClick={() => handleDelete(item.id)}>
                    <DeleteIcon />
                  </IconButton>
                }
                divider={index < knowledge.length - 1}
              >
                <ListItemText
                  primary={
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                      {getContentTypeIcon(item.content_type)}
                      {item.title}
                    </Box>
                  }
                  secondary={item.content.substring(0, 100) + (item.content.length > 100 ? '...' : '')}
                />
              </ListItem>
            ))}
          </List>
        )}
      </Paper>

      <Dialog open={open} onClose={closeDialog} maxWidth="md" fullWidth>
        <DialogTitle>Add Knowledge</DialogTitle>
        <DialogContent>
          <Tabs value={tabValue} onChange={(e, v) => setTabValue(v)} sx={{ mb: 2 }}>
            <Tab icon={<TextFieldsIcon />} label="Text" />
            <Tab icon={<CloudUploadIcon />} label="Upload File" />
            <Tab icon={<LinkIcon />} label="Website URL" />
          </Tabs>

          {/* Text Tab */}
          {tabValue === 0 && (
            <Box>
              <TextField
                fullWidth
                label="Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                margin="normal"
                required
              />
              <TextField
                fullWidth
                label="Content"
                value={content}
                onChange={(e) => setContent(e.target.value)}
                margin="normal"
                multiline
                rows={10}
                required
                helperText="Enter the information you want the chatbot to learn"
              />
            </Box>
          )}

          {/* File Upload Tab */}
          {tabValue === 1 && (
            <Box>
              <TextField
                fullWidth
                label="Title (optional)"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                margin="normal"
                helperText="Leave empty to use filename"
              />
              <Box sx={{ mt: 2 }}>
                <Button
                  variant="outlined"
                  component="label"
                  fullWidth
                  startIcon={<CloudUploadIcon />}
                >
                  {file ? file.name : 'Choose File (PDF, DOCX, TXT)'}
                  <input
                    type="file"
                    hidden
                    accept=".pdf,.docx,.txt"
                    onChange={handleFileChange}
                  />
                </Button>
                <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                  Supported formats: PDF, Word (DOCX), Text (TXT)
                </Typography>
              </Box>
            </Box>
          )}

          {/* URL Tab */}
          {tabValue === 2 && (
            <Box>
              <TextField
                fullWidth
                label="Website URL"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                margin="normal"
                required
                placeholder="https://example.com"
                helperText="Enter a URL to scrape content from"
              />
              <TextField
                fullWidth
                label="Title (optional)"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                margin="normal"
                helperText="Leave empty to use URL as title"
              />
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={closeDialog}>Cancel</Button>
          <Button onClick={handleSubmit} variant="contained" disabled={loading || !canSubmit()}>
            {loading ? 'Processing...' : 'Add'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default KnowledgeBase;
