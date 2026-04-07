import { useState, useEffect, useRef } from "react";
import { Link } from "react-router-dom";
import {
  Trash2,
  ShoppingBag,
  ArrowRight,
  Plus,
  Minus,
  ShoppingCart,
} from "lucide-react";

export default function Cart({ onCartUpdate, onCheckout }) {
  const [items, setItems] = useState([]);
  const [summary, setSummary] = useState(null);
  const [loading, setLoading] = useState(true);
  const [checkoutLoading, setCheckoutLoading] = useState(false);
  const [checkoutMessage, setCheckoutMessage] = useState("");
  const hasFetchedOnMount = useRef(false);

  useEffect(() => {
    if (hasFetchedOnMount.current) return;
    hasFetchedOnMount.current = true;
    fetchCart();
  }, []);

  const fetchCart = async () => {
    setLoading(true);
    try {
      const res = await fetch("http://localhost:3001/api/cart/details", {
        credentials: "include",
      });
      const data = await res.json();

      const cartData = data?.items;
      const summaryData = data?.summary;

      setItems(Array.isArray(cartData) ? cartData : []);
      setSummary(
        summaryData || {
          subtotal: 0,
          kdv_rate: 0,
          kdv: 0,
          total: 0,
        },
      );
    } catch {
      setItems([]);
      setSummary(null);
    }
    setLoading(false);
  };

  const updateQuantity = async (id, newQty) => {
    if (newQty < 1) {
      removeItem(id);
      return;
    }
    await fetch(`http://localhost:3001/api/cart/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ quantity: newQty }),
    });
    fetchCart();
    onCartUpdate();
  };

  const removeItem = async (id) => {
    await fetch(`http://localhost:3001/api/cart/${id}`, {
      method: "DELETE",
      credentials: "include",
    });
    fetchCart();
    onCartUpdate();
  };

  const clearCart = async () => {
    await fetch("http://localhost:3001/api/cart", {
      method: "DELETE",
      credentials: "include",
    });
    fetchCart();
    onCartUpdate();
  };

  const handleCheckout = async () => {
    setCheckoutLoading(true);
    setCheckoutMessage("");

    try {
      const res = await fetch("http://localhost:3001/api/cart/checkout", {
        method: "POST",
        credentials: "include",
      });
      const contentType = res.headers.get("content-type") || "";
      let data = null;

      if (contentType.includes("application/json")) {
        data = await res.json();
      } else {
        const text = await res.text();
        data = { error: text };
      }

      if (!res.ok) {
        setCheckoutMessage(data?.error || "Siparis olusturulamadi");
        return;
      }

      setCheckoutMessage(data?.message || "Siparis olusturuldu");
      await fetchCart();
      onCartUpdate();
    } catch {
      setCheckoutMessage("Sunucuya ulasilamadi");
    } finally {
      setCheckoutLoading(false);
    }
  };

  const subtotal = summary?.subtotal || 0;
  const kdvTotal = summary?.kdv || 0;
  const grandTotal = summary?.total || 0;
  const kdvRate = summary?.kdv_rate || 0;
  const kdvLines = Array.isArray(summary?.lines) ? summary.lines : [];
  const totalItems = items.reduce((sum, item) => sum + item.quantity, 0);

  if (loading) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-8">
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="h-24 rounded bg-gray-200 animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-8 text-center">
        <div className="py-16">
          {checkoutMessage && (
            <div className="mb-4 rounded-md border border-green-200 bg-green-50 px-4 py-2 text-green-700 text-sm inline-block">
              {checkoutMessage}
            </div>
          )}
          <div className="w-20 h-20 rounded-full bg-gray-200 flex items-center justify-center mx-auto mb-5">
            <ShoppingCart size={32} className="text-gray-400" />
          </div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">
            Sepetiniz Boş
          </h2>
          <p className="text-gray-500 mb-6">Hadi bir şeyler ekleyelim.</p>
          <Link
            to="/"
            className="inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md"
          >
            <ShoppingBag size={18} />
            Alışverişe Başla
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-5xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-5">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">Sepetim</h1>
          <p className="text-gray-500 text-sm mt-1">{totalItems} ürün</p>
        </div>
        <button
          onClick={clearCart}
          className="flex items-center gap-2 text-red-600 text-sm hover:underline"
        >
          <Trash2 size={14} />
          Sepeti Temizle
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div className="lg:col-span-2 space-y-3">
          {items.map((item) => (
            <div
              key={item.ID}
              className="flex gap-4 bg-white border border-gray-200 rounded-md p-4"
            >
              <Link to={`/product/${item.product_id}`} className="shrink-0">
                <div className="w-20 h-20 rounded-md overflow-hidden bg-gray-100">
                  <img
                    src={item.product.image}
                    alt={item.product.name}
                    className="w-full h-full object-cover"
                  />
                </div>
              </Link>

              <div className="flex-1 min-w-0">
                <Link to={`/product/${item.product_id}`}>
                  <h3 className="text-gray-800 font-semibold text-sm hover:underline">
                    {item.product.name}
                  </h3>
                </Link>
                <p className="text-gray-500 text-xs mt-0.5">
                  {item.product.category}
                </p>

                <div className="flex items-center justify-between mt-3">
                  <div className="flex items-center border border-gray-300 rounded-md overflow-hidden">
                    <button
                      onClick={() => updateQuantity(item.ID, item.quantity - 1)}
                      className="px-3 py-1 text-gray-700 hover:bg-gray-100"
                    >
                      <Minus size={13} />
                    </button>
                    <span className="px-3 py-1 text-gray-800 text-sm font-medium min-w-[2.5rem] text-center">
                      {item.quantity}
                    </span>
                    <button
                      onClick={() => updateQuantity(item.ID, item.quantity + 1)}
                      className="px-3 py-1 text-gray-700 hover:bg-gray-100"
                    >
                      <Plus size={13} />
                    </button>
                  </div>

                  <div className="text-right">
                    <p className="text-blue-700 font-bold text-sm">
                      ₺
                      {(item.product.price * item.quantity).toLocaleString(
                        "tr-TR",
                        { minimumFractionDigits: 2 },
                      )}
                    </p>
                    <p className="text-gray-500 text-xs">
                      ₺
                      {item.product.price.toLocaleString("tr-TR", {
                        minimumFractionDigits: 2,
                      })}{" "}
                      × {item.quantity}
                    </p>
                  </div>

                  <button
                    onClick={() => removeItem(item.ID)}
                    className="ml-3 text-gray-400 hover:text-red-600"
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="lg:col-span-1">
          <div className="bg-white border border-gray-200 rounded-md p-4 sticky top-20">
            <h2 className="text-gray-800 font-semibold mb-4">Sipariş Özeti</h2>

            <div className="space-y-3 mb-4">
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">
                  Ara Toplam ({totalItems} ürün)
                </span>
                <span className="text-gray-800">
                  ₺
                  {subtotal.toLocaleString("tr-TR", {
                    minimumFractionDigits: 2,
                  })}
                </span>
              </div>

              <div className="bg-gray-50 rounded-md p-3 space-y-2">
                <p className="text-gray-500 text-xs uppercase tracking-wide mb-1">
                  KDV Detayı (%{Math.round(kdvRate * 100)})
                </p>
                {kdvLines.map((line) => (
                  <div
                    key={line.cart_item_id}
                    className="flex justify-between text-xs text-gray-600"
                  >
                    <span className="truncate max-w-[140px]">
                      {line.product}
                    </span>
                    <span>
                      ₺
                      {line.kdv.toLocaleString("tr-TR", {
                        minimumFractionDigits: 2,
                      })}
                    </span>
                  </div>
                ))}
                <div className="border-t border-gray-200 pt-2 flex justify-between text-sm text-gray-800 font-medium">
                  <span>Toplam KDV</span>
                  <span>
                    ₺
                    {kdvTotal.toLocaleString("tr-TR", {
                      minimumFractionDigits: 2,
                    })}
                  </span>
                </div>
              </div>

              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Kargo</span>
                <span className="text-green-700 font-medium">Ücretsiz</span>
              </div>
            </div>

            <div className="border-t border-gray-200 pt-3 mb-4">
              <div className="flex justify-between items-baseline">
                <span className="text-gray-800 font-semibold">
                  Genel Toplam
                </span>
                <span className="text-xl font-bold text-gray-900">
                  ₺
                  {grandTotal.toLocaleString("tr-TR", {
                    minimumFractionDigits: 2,
                  })}
                </span>
              </div>
              <p className="text-gray-500 text-xs mt-1">KDV dahil</p>
            </div>

            <button
              onClick={handleCheckout}
              disabled={checkoutLoading}
              className="w-full flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-700 text-white py-2.5 rounded-md text-sm font-medium"
            >
              {checkoutLoading ? "Olusturuluyor..." : "Odemeye Gec"}
              <ArrowRight size={18} />
            </button>
            {checkoutMessage && (
              <p className="mt-2 text-center text-sm text-green-700">
                {checkoutMessage}
              </p>
            )}

            <Link
              to="/"
              className="flex items-center justify-center gap-2 mt-3 text-gray-500 hover:text-gray-700 text-sm"
            >
              <ShoppingBag size={14} />
              Alışverişe Devam Et
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
