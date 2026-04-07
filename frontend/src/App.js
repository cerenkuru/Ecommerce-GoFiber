import { useState, useEffect } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Navbar from "./Navbar";
import ProductList from "./ProductList";
import ProductDetail from "./ProductDetail";
import Cart from "./Cart";

export default function App() {
  const [cartCount, setCartCount] = useState(0);

  const incrementCartCount = (amount = 1) => {
    setCartCount((prev) => prev + amount);
  };

  const refreshCartCount = async () => {
    try {
      const res = await fetch("http://localhost:3001/api/cart", {
        credentials: "include",
      });
      const data = await res.json();
      const total = data.reduce((sum, item) => sum + item.quantity, 0);
      setCartCount(total);
    } catch {}
  };

  useEffect(() => {
    refreshCartCount();
  }, []);

  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-100">
        <Navbar cartCount={cartCount} />
        <Routes>
          <Route
            path="/"
            element={<ProductList onCartAdded={incrementCartCount} />}
          />
          <Route
            path="/product/:id"
            element={<ProductDetail onCartAdded={incrementCartCount} />}
          />
          <Route
            path="/cart"
            element={<Cart onCartUpdate={refreshCartCount} />}
          />
        </Routes>
      </div>
    </BrowserRouter>
  );
}
