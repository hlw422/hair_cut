import React, { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
} from 'recharts';
import { useQuery } from '@tanstack/react-query';
import { analyticsAPI } from '@/services/api';

// Mock数据（实际应从API获取）
const revenueData = [
  { date: '01/14', revenue: 42500, orders: 156 },
  { date: '01/15', revenue: 52300, orders: 189 },
  { date: '01/16', revenue: 48100, orders: 172 },
  { date: '01/17', revenue: 61200, orders: 215 },
  { date: '01/18', revenue: 58900, orders: 201 },
  { date: '01/19', revenue: 72400, orders: 268 },
  { date: '01/20', revenue: 68500, orders: 245 },
];

const orderSourceData = [
  { name: '微信小程序', value: 45, color: '#C8A882' },
  { name: '官网预约', value: 20, color: '#B8956A' },
  { name: '到店消费', value: 25, color: '#D4B896' },
  { name: '电话预订', value: 10, color: '#E8D5B8' },
];

const storeRankData = [
  { name: '静安店', gmv: 128500 },
  { name: '徐汇店', gmv: 115200 },
  { name: '浦东店', gmv: 98700 },
  { name: '虹桥店', gmv: 87600 },
  { name: '杨浦店', gmv: 76800 },
];

const recentOrders = [
  { id: 'ORD001', user: '王**', store: '静安店', amount: 198.0, status: '已完成', time: '2分钟前' },
  { id: 'ORD002', user: '李**', store: '徐汇店', amount: 388.0, status: '进行中', time: '5分钟前' },
  { id: 'ORD003', user: '张**', store: '浦东店', amount: 168.0, status: '待支付', time: '12分钟前' },
  { id: 'ORD004', user: '刘**', store: '静安店', amount: 528.0, status: '已完成', time: '18分钟前' },
  { id: 'ORD005', user: '陈**', store: '虹桥店', amount: 98.0, status: '已退款', time: '25分钟前' },
];

interface KpiCardProps {
  title: string;
  value: string | number;
  change: number;
  changeLabel: string;
  icon: React.ReactNode;
  color?: string;
}

function KpiCard({ title, value, change, changeLabel, icon, color = '#C8A882' }: KpiCardProps) {
  const isPositive = change >= 0;

  return (
    <Card className="border-0 shadow-card hover:shadow-card-hover transition-shadow duration-300">
      <CardContent className="p-6">
        <div className="flex items-start justify-between">
          <div className="space-y-2">
            <p className="text-sm font-medium text-[#6B7280]">{title}</p>
            <p className="text-3xl font-bold text-[#1A1A1A]">{value}</p>
            <div className="flex items-center gap-2 text-sm">
              <span
                className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
                  isPositive
                    ? 'bg-green-50 text-green-700'
                    : 'bg-red-50 text-red-700'
                }`}
              >
                {isPositive ? '↑' : '↓'} {Math.abs(change)}%
              </span>
              <span className="text-[#9CA3AF]">{changeLabel}</span>
            </div>
          </div>
          <div
            className="w-14 h-14 rounded-2xl flex items-center justify-center"
            style={{ backgroundColor: `${color}15` }}
          >
            <div style={{ color }}>{icon}</div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export default function DashboardPage() {
  const [dateRange, setDateRange] = useState('today');

  // TODO: 使用TanStack Query获取真实数据
  // const { data: dashboardData } = useQuery({
  //   queryKey: ['dashboard-stats'],
  //   queryFn: () => analyticsAPI.getDashboardStats(),
  // });

  return (
    <div className="space-y-6 p-6">
      {/* 页面标题 */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-[#1A1A1A]">运营概览</h1>
          <p className="text-[#6B7280] mt-1">欢迎回来！以下是今日业务数据总览</p>
        </div>
        <div className="flex gap-2">
          {/* 时间范围选择器 */}
          {['today', 'week', 'month'].map((range) => (
            <button
              key={range}
              onClick={() => setDateRange(range)}
              className={`px-4 py-2 rounded-xl text-sm font-medium transition-all ${
                dateRange === range
                  ? 'bg-[#C8A882] text-white shadow-md'
                  : 'bg-white text-[#6B7280] hover:bg-gray-50 border border-gray-200'
              }`}
            >
              {{ today: '今日', week: '本周', month: '本月' }[range]}
            </button>
          ))}
        </div>
      </div>

      {/* KPI指标卡片区 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <KpiCard
          title="今日GMV (元)"
          value={168,520}
          change={12.5}
          changeLabel="较昨日"
          icon={<span className="text-2xl">💰</span>}
          color="#10B981"
        />
        <KpiCard
          title="订单量 (笔)"
          value={586}
          change={8.3}
          changeLabel="较昨日"
          icon={<span className="text-2xl">📦</span>}
          color="#3B82F6"
        />
        <KpiCard
          title="活跃用户"
          value={1,280}
          change={15.2}
          changeLabel="较昨日"
          icon={<span className="text-2xl">👥</span>}
          color="#8B5CF6"
        />
        <KpiCard
          title="转化率 (%)"
          value={18.6}
          change={-2.1}
          changeLabel="较昨日"
          icon={<span className="text-2xl">📈</span>}
          color="#F59E0B"
        />
      </div>

      {/* 中部图表区 */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* 营收趋势图 - 占2列 */}
        <Card className="lg:col-span-2 border-0 shadow-card">
          <CardHeader className="pb-2">
            <CardTitle className="text-lg font-semibold">营收趋势</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={320}>
              <LineChart data={revenueData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis dataKey="date" stroke="#9CA3AF" fontSize={12} />
                <YAxis stroke="#9CA3AF" fontSize={12} />
                <Tooltip
                  contentStyle={{
                    borderRadius: 12,
                    boxShadow: '0 4px 16px rgba(0,0,0,0.1)',
                    border: 'none',
                  }}
                />
                <Line
                  type="monotone"
                  dataKey="revenue"
                  stroke="#C8A882"
                  strokeWidth={3}
                  dot={{ fill: '#C8A882', r: 4 }}
                  activeDot={{ r: 6, stroke: '#C8A882', strokeWidth: 2 }}
                />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* 订单来源饼图 */}
        <Card className="border-0 shadow-card">
          <CardHeader className="pb-2">
            <CardTitle className="text-lg font-semibold">订单来源分布</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={320}>
              <PieChart>
                <Pie
                  data={orderSourceData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={100}
                  paddingAngle={5}
                  dataKey="value"
                >
                  {orderSourceData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
            <div className="space-y-2 mt-4">
              {orderSourceData.map((item, index) => (
                <div key={index} className="flex items-center justify-between text-sm">
                  <div className="flex items-center gap-2">
                    <div
                      className="w-3 h-3 rounded-full"
                      style={{ backgroundColor: item.color }}
                    />
                    <span className="text-[#374151]">{item.name}</span>
                  </div>
                  <span className="font-medium text-[#1A1A1A]">{item.value}%</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* 底部区域 */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* 门店排行榜 */}
        <Card className="lg:col-span-2 border-0 shadow-card">
          <CardHeader className="pb-2">
            <CardTitle className="text-lg font-semibold">门店GMV排行 TOP5</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={280}>
              <BarChart data={storeRankData} layout="vertical">
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis type="number" stroke="#9CA3AF" fontSize={12} tickFormatter={(v) => `${(v / 10000).toFixed(1)}万`} />
                <YAxis type="category" dataKey="name" stroke="#9CA3AF" fontSize={12} width={60} />
                <Tooltip
                  formatter={(value: number) => [`¥${value.toLocaleString()}`, 'GMV']}
                  contentStyle={{
                    borderRadius: 12,
                    boxShadow: '0 4px 16px rgba(0,0,0,0.1)',
                    border: 'none',
                  }}
                />
                <Bar dataKey="gmv" fill="#C8A882" radius={[0, 6, 6, 0]} barSize={24} />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* 实时订单动态 */}
        <Card className="border-0 shadow-card">
          <CardHeader className="pb-2 flex flex-row items-center justify-between">
            <CardTitle className="text-lg font-semibold">实时订单</CardTitle>
            <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
          </CardHeader>
          <CardContent className="p-0">
            <div className="divide-y divide-[#F3F4F6]">
              {recentOrders.map((order) => (
                <div key={order.id} className="p-4 hover:bg-gray-50 transition-colors">
                  <div className="flex items-center justify-between mb-1">
                    <span className="font-mono text-xs text-[#9CA3AF]">{order.id}</span>
                    <span className="text-xs text-[#9CA3AF]">{order.time}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-[#374151]">{order.user} · {order.store}</span>
                    <div className="flex items-center gap-2">
                      <span className="font-semibold text-[#1A1A1A]">¥{order.amount}</span>
                      <span
                        className={`text-xs px-2 py-0.5 rounded ${
                          order.status === '已完成'
                            ? 'bg-green-50 text-green-600'
                            : order.status === '进行中'
                            ? 'bg-blue-50 text-blue-600'
                            : order.status === '待支付'
                            ? 'bg-yellow-50 text-yellow-600'
                            : 'bg-red-50 text-red-600'
                        }`}
                      >
                        {order.status}
                      </span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            <button className="w-full p-3 text-sm text-[#C8A882] hover:bg-[#FBF7F0] transition-colors rounded-b-xl">
              查看全部订单 →
            </button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
