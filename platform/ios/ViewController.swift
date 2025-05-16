import UIKit
import WebKit

class ViewController: UIViewController, WKNavigationDelegate, WKScriptMessageHandler {
    private var webView: WKWebView!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        // Configure WKWebView
        let configuration = WKWebViewConfiguration()
        let userContentController = WKUserContentController()
        
        // Add script message handler for JS bridge
        userContentController.add(self, name: "iOSBridge")
        configuration.userContentController = userContentController
        
        // Enable developer tools in debug mode
        #if DEBUG
        if #available(iOS 16.4, *) {
            configuration.preferences.isInspectable = true
        }
        #endif
        
        // Create WKWebView
        webView = WKWebView(frame: view.bounds, configuration: configuration)
        webView.navigationDelegate = self
        webView.autoresizingMask = [.flexibleWidth, .flexibleHeight]
        view.addSubview(webView)
        
        // Load web app
        loadWebApp()
    }
    
    private func loadWebApp() {
        // Check for development mode
        #if DEBUG
        // For development, load from dev server
        if let url = URL(string: "http://localhost:3001") {
            let request = URLRequest(url: url)
            webView.load(request)
        }
        #else
        // For production, load from bundled files
        if let htmlPath = Bundle.main.path(forResource: "index", ofType: "html", inDirectory: "assets") {
            let htmlUrl = URL(fileURLWithPath: htmlPath)
            webView.loadFileURL(htmlUrl, allowingReadAccessTo: htmlUrl.deletingLastPathComponent())
        }
        #endif
    }
    
    // MARK: - WKScriptMessageHandler
    
    func userContentController(_ userContentController: WKUserContentController, didReceive message: WKScriptMessage) {
        guard let dict = message.body as? [String: Any] else { return }
        
        if message.name == "iOSBridge", let action = dict["action"] as? String {
            if action == "getPlatformInfo" {
                let deviceInfo = "iOS \(UIDevice.current.systemVersion)"
                // Call JavaScript function to set platform info
                webView.evaluateJavaScript("setPlatformInfo('\(deviceInfo)')", completionHandler: nil)
            }
        }
    }
    
    // MARK: - WKNavigationDelegate
    
    func webView(_ webView: WKWebView, decidePolicyFor navigationAction: WKNavigationAction, decisionHandler: @escaping (WKNavigationActionPolicy) -> Void) {
        if navigationAction.navigationType == .linkActivated {
            // Handle external links
            if let url = navigationAction.request.url, UIApplication.shared.canOpenURL(url) {
                UIApplication.shared.open(url)
                decisionHandler(.cancel)
                return
            }
        }
        decisionHandler(.allow)
    }
} 