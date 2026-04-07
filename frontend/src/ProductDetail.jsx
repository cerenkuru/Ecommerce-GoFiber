import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { ShoppingCart, ArrowLeft, Package, CheckCircle } from "lucide-react";

export default function ProductDetail({ onCartAdded }) {
  const { id } = useParams();
  const navigate = useNavigate();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [quantity, setQuantity] = useState(1);
  const [added, setAdded] = useState(false);

  useEffect(() => {
    fetch(`http://localhost:3001/api/products/${id}`)
      .then((r) => r.json())
      .then((data) => {
        setProduct(data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [id]);

  const addToCart = async () => {
    try {
      await fetch("http://localhost:3001/api/cart", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ product_id: product.ID, quantity }),
      });
      setAdded(true);
      setTimeout(() => setAdded(false), 2000);
      onCartAdded(quantity);
    } catch (err) {
      console.error(err);
    }
  };

  if (loading) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-10">
          <div className="aspect-square rounded-md bg-gray-200 animate-pulse" />
          <div className="space-y-4">
            {[...Array(5)].map((_, i) => (
              <div
                key={i}
                className="h-8 rounded bg-gray-200 animate-pulse"
                style={{ width: `${80 - i * 10}%` }}
              />
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="text-center py-20 text-gray-500">
        <Package size={48} className="mx-auto mb-4" />
        <p>Ürün bulunamadı</p>
      </div>
    );
  }

  const kdv = product.price * quantity * 0.2;
  const total = product.price * quantity + kdv;

  return (
    <div className="max-w-5xl mx-auto px-4 py-8">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-6 text-sm"
      >
        <ArrowLeft size={16} />
        Geri Dön
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div>
          <div className="aspect-square rounded-md overflow-hidden bg-white border border-gray-200">
            <img
              src={product.image}
              alt={product.name}
              className="w-full h-full object-cover"
            />
          </div>
        </div>

        <div>
          <div className="flex items-center gap-3 mb-3">
            <span className="bg-gray-100 text-gray-600 text-xs px-2 py-1 rounded">
              {product.category}
            </span>
          </div>

          <h1 className="text-2xl font-bold text-gray-800 mb-3">
            {product.name}
          </h1>
          <p className="text-gray-600 mb-6">{product.description}</p>

          <div className="bg-white border border-gray-200 rounded-md p-4 mb-5">
            <div className="flex items-baseline gap-2">
              <span className="text-2xl font-bold text-blue-700">
                ₺
                {product.price.toLocaleString("tr-TR", {
                  minimumFractionDigits: 2,
                })}
              </span>
              <span className="text-sm text-gray-500">/ adet</span>
            </div>
            <p className="text-sm text-gray-500 mt-1">
              Stokta: {product.stock} adet
            </p>
          </div>

          <div className="flex items-center gap-4 mb-5">
            <span className="text-sm text-gray-600">Adet:</span>
            <div className="flex items-center border border-gray-300 rounded-md overflow-hidden bg-white">
              <button
                onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                className="px-3 py-1 text-gray-700 hover:bg-gray-100"
              >
                −
              </button>
              <span className="px-4 py-1 min-w-[3rem] text-center text-gray-800 font-medium">
                {quantity}
              </span>
              <button
                onClick={() =>
                  setQuantity((q) => Math.min(product.stock, q + 1))
                }
                className="px-3 py-1 text-gray-700 hover:bg-gray-100"
              >
                +
              </button>
            </div>
          </div>

          {quantity > 1 && (
            <div className="bg-white border border-gray-200 rounded-md px-4 py-3 mb-5 space-y-1 text-sm text-gray-600">
              <div className="flex justify-between">
                <span>Ara Toplam ({quantity} adet)</span>
                <span>
                  ₺
                  {(product.price * quantity).toLocaleString("tr-TR", {
                    minimumFractionDigits: 2,
                  })}
                </span>
              </div>
              <div className="flex justify-between">
                <span>KDV (%20)</span>
                <span>
                  ₺{kdv.toLocaleString("tr-TR", { minimumFractionDigits: 2 })}
                </span>
              </div>
              <div className="flex justify-between font-semibold text-gray-800 border-t border-gray-200 pt-1 mt-1">
                <span>Toplam</span>
                <span>
                  ₺{total.toLocaleString("tr-TR", { minimumFractionDigits: 2 })}
                </span>
              </div>
            </div>
          )}

          <button
            onClick={addToCart}
            disabled={product.stock === 0}
            className={`w-full flex items-center justify-center gap-2 py-2.5 rounded-md text-sm font-medium ${
              added
                ? "bg-green-100 text-green-700"
                : product.stock === 0
                  ? "bg-gray-200 text-gray-500 cursor-not-allowed"
                  : "bg-blue-600 hover:bg-blue-700 text-white"
            }`}
          >
            {added ? (
              <>
                <CheckCircle size={20} />
                Sepete Eklendi!
              </>
            ) : (
              <>
                <ShoppingCart size={20} />
                {product.stock === 0 ? "Stokta Yok" : "Sepete Ekle"}
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
