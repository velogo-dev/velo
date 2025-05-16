package com.example.golangmobile

import android.annotation.SuppressLint
import android.os.Build
import android.os.Bundle
import android.webkit.*
import androidx.appcompat.app.AppCompatActivity
import android.content.Context
import android.content.Intent
import android.net.Uri
import android.webkit.WebView
import android.webkit.WebViewClient
import android.widget.Toast

class MainActivity : AppCompatActivity() {
    private lateinit var webView: WebView

    @SuppressLint("SetJavaScriptEnabled")
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        webView = findViewById(R.id.webview)

        // Enable JavaScript and debugging
        webView.settings.javaScriptEnabled = true
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.KITKAT) {
            WebView.setWebContentsDebuggingEnabled(true)
        }

        // Configure cache and security
        webView.settings.domStorageEnabled = true
        webView.settings.allowFileAccess = true

        // Set up app-specific WebViewClient
        webView.webViewClient = object : WebViewClient() {
            override fun shouldOverrideUrlLoading(view: WebView?, request: WebResourceRequest?): Boolean {
                val url = request?.url.toString()
                // Handle external links if needed
                if (url.startsWith("http://") || url.startsWith("https://") && !url.contains("localhost")) {
                    val intent = Intent(Intent.ACTION_VIEW, Uri.parse(url))
                    startActivity(intent)
                    return true
                }
                return false
            }
        }

        // Add JavaScript interface for communication
        webView.addJavascriptInterface(WebAppInterface(this), "AndroidBridge")

        // Load the web app
        loadWebApp()
    }

    private fun loadWebApp() {
        // Check for development mode (could be set via build config)
        val isDevelopment = true // Change based on build type
        
        if (isDevelopment) {
            // For development, load from development server
            webView.loadUrl("http://localhost:3000") // Android emulator localhost equivalent
        } else {
            // For production, load from assets (the bundled web app)
            webView.loadUrl("file:///android_asset/index.html")
        }
    }

    // Handle back button in WebView
    override fun onBackPressed() {
        if (webView.canGoBack()) {
            webView.goBack()
        } else {
            super.onBackPressed()
        }
    }

    // JavaScript interface for communication between JS and Android
    private class WebAppInterface(private val context: Context) {
        @JavascriptInterface
        fun getPlatformInfo(): String {
            return "Android ${Build.VERSION.RELEASE} (SDK ${Build.VERSION.SDK_INT})"
        }

        @JavascriptInterface
        fun showToast(message: String) {
            Toast.makeText(context, message, Toast.LENGTH_SHORT).show()
        }
    }
} 