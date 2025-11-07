import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Box,
  Button,
} from '@mui/material';
import {
  Chat as ChatIcon,
  Add as AddIcon,
} from '@mui/icons-material';
import api from '../../services/api';

function Dashboard() {
  const [chatbots, setChatbots] = useState([]);
  const [loading, setLoading] = useState(true);
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
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
        <Typography variant="h4">Dashboard</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => navigate('/chatbots')}
        >
          Create Chatbot
        </Button>
      </Box>

      <Grid container spacing={3}>
        <Grid item xs={12} sm={6} md={4}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <ChatIcon sx={{ fontSize: 40, mr: 2, color: 'primary.main' }} />
                <Box>
                  <Typography variant="h4">{chatbots.length}</Typography>
                  <Typography variant="body2" color="text.secondary">
                    Total Chatbots
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Box sx={{ mt: 4 }}>
        <Typography variant="h5" gutterBottom>
          Recent Chatbots
        </Typography>
        {loading ? (
          <Typography>Loading...</Typography>
        ) : chatbots.length === 0 ? (
          <Typography color="text.secondary">
            No chatbots yet. Create your first chatbot to get started!
          </Typography>
        ) : (
          <Grid container spacing={2}>
            {chatbots.slice(0, 6).map((chatbot) => (
              <Grid item xs={12} sm={6} md={4} key={chatbot.id}>
                <Card sx={{ cursor: 'pointer' }} onClick={() => navigate(`/chatbots/${chatbot.id}`)}>
                  <CardContent>
                    <Typography variant="h6">{chatbot.name}</Typography>
                    <Typography variant="body2" color="text.secondary">
                      {chatbot.description || 'No description'}
                    </Typography>
                    <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                      Status: {chatbot.status}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        )}
      </Box>
    </Box>
  );
}

export default Dashboard;


