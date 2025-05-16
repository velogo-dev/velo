import { useState } from "react";
import "./App.css";

function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="app">
      <h1>Golang Mobile Framework</h1>
      <p>Build with Go + React</p>

      <div className="card">
        <button onClick={() => setCount(count + 1)}>Count is {count}</button>
        <p>
          Edit <code>src/App.jsx</code> and save to test hot reload
        </p>
      </div>

      <p className="platform-info" id="platform-info">
        Platform detection will appear here
      </p>
    </div>
  );
}

export default App;
