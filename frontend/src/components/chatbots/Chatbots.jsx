import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Card,
  CardContent,
  CardActions,
  Typography,
  Grid,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  IconButton,
} from '@mui/material';
import { Add as AddIcon, Delete as DeleteIcon, Edit as EditIcon } from '@mui/icons-material';
import api from '../../services/api';

function Chatbots() {
  const [chatbots, setChatbots] = useState([]);
  const [open, setOpen] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    fetchChatbots();
  }, []);

  const fetchChatbots = async () => {
    try {
      const response = await api.get('/chatbots');
      setChatbots(response.data || []);
    } catch (error) {
      console.error('Failed to fetch chatbots:', error);
    }
  };

  const handleCreate = async () => {
    setError('');
    setLoading(true);

    try {
      await api.post('/chatbots', { name, description });
      setOpen(false);
      setName('');
      setDescription('');
      fetchChatbots();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create chatbot');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this chatbot?')) {
      return;
    }

    try {
      await api.delete(`/chatbots/${id}`);
      fetchChatbots();
    } catch (error) {
      console.error('Failed to delete chatbot:', error);
    }
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
        <Typography variant="h4">Chatbots</Typography>
        <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpen(true)}>
          Create Chatbot
        </Button>
      </Box>

      <Grid container spacing={3}>
        {chatbots.map((chatbot) => (
          <Grid item xs={12} sm={6} md={4} key={chatbot.id}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  {chatbot.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {chatbot.description || 'No description'}
                </Typography>
                <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                  Status: {chatbot.status}
                </Typography>
              </CardContent>
              <CardActions>
                <Button size="small" onClick={() => navigate(`/chatbots/${chatbot.id}`)}>
                  Manage
                </Button>
                <Button size="small" onClick={() => navigate(`/knowledge/${chatbot.id}`)}>
                  Knowledge
                </Button>
                <IconButton size="small" color="error" onClick={() => handleDelete(chatbot.id)}>
                  <DeleteIcon />
                </IconButton>
              </CardActions>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Dialog open={open} onClose={() => setOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Create New Chatbot</DialogTitle>
        <DialogContent>
          {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
          <TextField
            fullWidth
            label="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            margin="normal"
            required
          />
          <TextField
            fullWidth
            label="Description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            margin="normal"
            multiline
            rows={3}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Cancel</Button>
          <Button onClick={handleCreate} variant="contained" disabled={loading || !name}>
            {loading ? 'Creating...' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default Chatbots;

