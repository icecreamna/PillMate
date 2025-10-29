import { useState } from "react";
import { logout as apiLogout } from "../services/auth";

export default function LogoutButton({
  className,
  children = "ออกจากระบบ",
  onAfter,
  confirm = false,
}) {
  const [loading, setLoading] = useState(false);

  const onClick = async () => {
    if (loading) return;
    if (confirm && !window.confirm("ยืนยันออกจากระบบ?")) return;

    setLoading(true);
    try { await apiLogout(); } catch {}
    localStorage.removeItem("role");
    localStorage.removeItem("email");
    if (typeof onAfter === "function") onAfter();
    else window.location.href = "/login";
    setLoading(false);
  };

  return (
    <button className={className} onClick={onClick} disabled={loading}>
      {loading ? "กำลังออกจากระบบ..." : children}
    </button>
  );
}
