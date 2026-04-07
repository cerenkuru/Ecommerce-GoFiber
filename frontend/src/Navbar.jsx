import { Link } from "react-router-dom";
import { ShoppingCart, Zap } from "lucide-react";

export default function Navbar({ cartCount }) {
  return (
    <nav className="sticky top-0 z-30 bg-white border-b border-gray-200">
      <div className="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between">
        <Link
          to="/"
          className="flex items-center gap-2 text-gray-800 font-semibold"
        >
          <div className="w-8 h-8 rounded-md bg-blue-600 flex items-center justify-center">
            <Zap size={16} className="text-white" />
          </div>
          <span>Cero Shop</span>
        </Link>

        <div className="hidden md:flex items-center gap-5 text-sm">
          {["Elektronik", "Bilgisayar", "Aksesuar", "Kamera"].map((cat) => (
            <Link
              key={cat}
              to={`/?category=${cat}`}
              className="text-gray-600 hover:text-gray-900"
            >
              {cat}
            </Link>
          ))}
        </div>

        <Link
          to="/cart"
          className="flex items-center gap-2 border border-gray-300 rounded-md px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50"
        >
          <ShoppingCart size={16} />
          <span>Sepet</span>
          <span className="inline-flex items-center justify-center min-w-5 h-5 px-1 rounded-full bg-blue-600 text-white text-xs">
            {cartCount > 9 ? "9+" : cartCount}
          </span>
        </Link>
      </div>
    </nav>
  );
}
