import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useAuthStore } from '@/store/auth';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuthStore();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!username || !password) {
      return;
    }

    setLoading(true);
    try {
      await login(username, password);
      navigate('/dashboard');
    } catch (error) {
      console.error('登录失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-[#F7F8FA] via-white to-[#FBF7F0]">
      {/* 背景装饰 */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-[500px] h-[500px] rounded-full bg-[#C8A882]/10 blur-3xl" />
        <div className="absolute -bottom-40 -left-40 w-[400px] h-[400px] rounded-full bg-[#C8A882]/5 blur-3xl" />
      </div>

      <Card className="w-full max-w-md mx-4 shadow-xl border-0 relative z-10">
        {/* 品牌Logo区 */}
        <CardHeader className="text-center pb-2">
          <div className="flex items-center justify-center mb-4">
            <div className="w-16 h-16 bg-gradient-to-br from-[#C8A882] to-[#B8956A] rounded-2xl flex items-center justify-center shadow-lg shadow-[#C8A882]/20">
              <span className="text-white text-3xl font-bold">H</span>
            </div>
          </div>
          <CardTitle className="text-2xl font-bold text-[#1A1A1A]">HairCut 运营后台</CardTitle>
          <CardDescription className="text-[#6B7280] mt-1">
            连锁理发店数字化管理平台
          </CardDescription>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="space-y-2">
              <label className="text-sm font-medium text-[#374151]">账号</label>
              <Input
                type="text"
                placeholder="请输入管理员账号"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="h-12 border-[#E5E7EB] focus:border-[#C8A882] focus:ring-[#C8A882]/20 rounded-xl"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-[#374151]">密码</label>
              <Input
                type="password"
                placeholder="请输入密码"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="h-12 border-[#E5E7EB] focus:border-[#C8A882] focus:ring-[#C8A882]/20 rounded-xl"
                onKeyPress={(e) => e.key === 'Enter' && handleSubmit(e)}
              />
            </div>

            <div className="flex items-center justify-between text-sm">
              <label className="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" className="rounded border-gray-300 text-[#C8A882] focus:ring-[#C8A882]" />
                <span className="text-[#6B7280]">记住登录状态</span>
              </label>
              <button type="button" className="text-[#C8A882] hover:text-[#B8956A] font-medium transition-colors">
                忘记密码？
              </button>
            </div>

            <Button
              type="submit"
              disabled={loading || !username || !password}
              className="w-full h-12 bg-gradient-to-r from-[#C8A882] to-[#B8956A] hover:from-[#B8956A] hover:to-[#94784D] text-white font-semibold rounded-xl shadow-lg shadow-[#C8A882]/25 hover:shadow-xl hover:shadow-[#C8A882]/30 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? (
                <span className="flex items-center justify-center gap-2">
                  <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                  </svg>
                  登录中...
                </span>
              ) : '登 录'}
            </Button>
          </form>

          {/* 底部提示 */}
          <div className="mt-6 pt-6 border-t border-[#F3F4F6] text-center text-sm text-[#9CA3AF]">
            <p>© 2024 HairCut Team · 企业版 v1.0.0</p>
            <p className="mt-1">如有问题请联系技术支持</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
