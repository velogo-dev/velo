import React, { useState, useEffect } from 'react';
import { StyleSheet, Text, View, TouchableOpacity, Platform } from 'react-native';
import { StatusBar } from 'expo-status-bar';
import * as WebBrowser from 'expo-web-browser';

export default function App() {
      const [count, setCount] = useState(0);
      const [platformInfo, setPlatformInfo] = useState('');

      useEffect(() => {
            // Detect platform
            setPlatformInfo(`Running on ${Platform.OS} (${Platform.Version})`);
      }, []);

      const openWebView = async () => {
            // Open the web version in WebBrowser
            await WebBrowser.openBrowserAsync('http://localhost:3001');
      };

      return (
            <View style={styles.container}>
                  <Text style={styles.title}>Golang Mobile Framework</Text>
                  <Text style={styles.subtitle}>Build with Go + React Native + Expo</Text>

                  <View style={styles.card}>
                        <TouchableOpacity style={styles.button} onPress={() => setCount(count + 1)}>
                              <Text style={styles.buttonText}>Count is {count}</Text>
                        </TouchableOpacity>
                        <Text style={styles.info}>Edit App.js and save to test hot reload</Text>
                  </View>

                  <Text style={styles.platformInfo}>{platformInfo}</Text>

                  <TouchableOpacity style={styles.webButton} onPress={openWebView}>
                        <Text style={styles.buttonText}>Open Web Version</Text>
                  </TouchableOpacity>

                  <StatusBar style="auto" />
            </View>
      );
}

const styles = StyleSheet.create({
      container: {
            flex: 1,
            backgroundColor: '#242424',
            alignItems: 'center',
            justifyContent: 'center',
            padding: 20,
      },
      title: {
            fontSize: 24,
            fontWeight: 'bold',
            marginBottom: 8,
            color: '#ffffff',
      },
      subtitle: {
            fontSize: 16,
            marginBottom: 24,
            color: '#ffffff',
      },
      card: {
            backgroundColor: '#1e1e1e',
            borderRadius: 8,
            padding: 20,
            width: '100%',
            maxWidth: 400,
            alignItems: 'center',
            marginVertical: 20,
      },
      button: {
            backgroundColor: '#646cff',
            paddingHorizontal: 20,
            paddingVertical: 10,
            borderRadius: 8,
            marginBottom: 16,
      },
      webButton: {
            backgroundColor: '#28a745',
            paddingHorizontal: 20,
            paddingVertical: 10,
            borderRadius: 8,
            marginTop: 20,
      },
      buttonText: {
            color: 'white',
            fontSize: 16,
            fontWeight: '500',
      },
      info: {
            color: '#cccccc',
            fontSize: 14,
            textAlign: 'center',
      },
      platformInfo: {
            color: '#ffffff',
            fontWeight: 'bold',
            marginTop: 24,
      },
}); 