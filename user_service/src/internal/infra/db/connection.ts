import mongoose from 'mongoose';
import config from '../../config/index.js';

let dbInstance: mongoose.Connection | null = null;

export async function connectDB(): Promise<mongoose.Connection> {
  if (dbInstance) {
    return dbInstance;
  }
  await mongoose.connect(config.db_url as string);
  dbInstance = mongoose.connection;
  dbInstance.on('connected', () => {
    console.log('MongoDB connected');
  });

  dbInstance.on('error', err => {
    console.error('MongoDB connection error:', err);
  });
  return dbInstance;
}
