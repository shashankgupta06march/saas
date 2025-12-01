import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Paper,
  Typography,
  List,
  ListItem,
  ListItemText,
  Button,
  Divider,
} from '@mui/material';
import { ArrowBack as ArrowBackIcon } from '@mui/icons-material';
import api from '../../services/api';

function Analytics() {
  const { chatbotId } = useParams();
  const navigate = useNavigate();
  const [conversations, setConversations] = useState([]);
  const [selectedConversation, setSelectedConversation] = useState(null);
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    fetchConversations();
  }, [chatbotId]);

  const fetchConversations = async () => {
    try {
      const response = await api.get(`/analytics/conversations/${chatbotId}`);
      setConversations(response.data || []);
    } catch (error) {
      console.error('Failed to fetch conversations:', error);
    }
  };

  const fetchMessages = async (conversationId) => {
    try {
      const response = await api.get(`/analytics/messages/${conversationId}`);
      setMessages(response.data || []);
      setSelectedConversation(conversationId);
    } catch (error) {
      console.error('Failed to fetch messages:', error);
    }
  };

  return (
    <Box>
      <Button startIcon={<ArrowBackIcon />} onClick={() => navigate('/chatbots')} sx={{ mb: 2 }}>
        Back to Chatbots
      </Button>

      <Typography variant="h4" gutterBottom>
        Analytics & Conversations
      </Typography>

      <Box sx={{ display: 'flex', gap: 2, mt: 3 }}>
        <Paper sx={{ flex: 1, p: 2 }}>
          <Typography variant="h6" gutterBottom>
            Conversations
          </Typography>
          {conversations.length === 0 ? (
            <Typography color="text.secondary">No conversations yet</Typography>
          ) : (
            <List>
              {conversations.map((conv, index) => (
                <ListItem
                  key={conv.id}
                  button
                  selected={selectedConversation === conv.id}
                  onClick={() => fetchMessages(conv.id)}
                  divider={index < conversations.length - 1}
                >
                  <ListItemText
                    primary={`Session: ${conv.session_id.substring(0, 8)}...`}
                    secondary={new Date(conv.started_at).toLocaleString()}
                  />
                </ListItem>
              ))}
            </List>
          )}
        </Paper>

        <Paper sx={{ flex: 2, p: 2 }}>
          <Typography variant="h6" gutterBottom>
            Messages
          </Typography>
          {!selectedConversation ? (
            <Typography color="text.secondary">
              Select a conversation to view messages
            </Typography>
          ) : messages.length === 0 ? (
            <Typography color="text.secondary">No messages</Typography>
          ) : (
            <Box>
              {messages.map((msg, index) => (
                <Box key={msg.id} sx={{ mb: 2 }}>
                  <Typography
                    variant="subtitle2"
                    color={msg.role === 'user' ? 'primary' : 'secondary'}
                  >
                    {msg.role === 'user' ? 'User' : 'Assistant'}
                  </Typography>
                  <Paper sx={{ p: 2, bgcolor: msg.role === 'user' ? '#e3f2fd' : '#f5f5f5' }}>
                    <Typography variant="body2">{msg.content}</Typography>
                    <Typography variant="caption" color="text.secondary">
                      {new Date(msg.timestamp).toLocaleString()}
                    </Typography>
                  </Paper>
                </Box>
              ))}
            </Box>
          )}
        </Paper>
      </Box>
    </Box>
  );
}

export default Analytics;



