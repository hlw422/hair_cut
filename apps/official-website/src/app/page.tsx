import Image from 'next/image';
import Link from 'next/link';
import { MapPin, Phone, Clock, Star, ArrowRight, Users, Scissors, Sparkles } from 'lucide-react';

// Mock数据
const featuredStores = [
  {
    id: 1,
    name: 'HairCut 精品沙龙（静安旗舰店）',
    city: '上海',
    address: '静安区南京西路123号',
    image: '/images/store-1.jpg',
    rating: 4.9,
    reviews: 520,
    avgPrice: 198,
    tags: ['环境优雅', '明星理发师', '免费停车'],
  },
  {
    id: 2,
    name: 'HairCut 造型中心（徐汇店）',
    city: '上海',
    address: '徐汇区淮海中路456号',
    image: '/images/store-2.jpg',
    rating: 4.9,
    reviews: 389,
    avgPrice: 178,
    tags: ['人气TOP1', '地铁直达'],
  },
  {
    id: 3,
    name: 'HairCut 概念店（浦东陆家嘴）',
    city: '上海',
    address: '浦东新区世纪大道789号',
    image: '/images/store-3.jpg',
    rating: 4.8,
    reviews: 267,
    avgPrice: 168,
    tags: ['江景沙龙', '新店开业优惠'],
  },
];

const topStylists = [
  {
    id: 1,
    name: 'Kevin Chen',
    title: '艺术总监 · 12年经验',
    specialty: ['日韩风格', '色彩设计', '新娘造型'],
    avatar: '/images/stylist-kevin.jpg',
    rating: 4.9,
    fansCount: '2.8k',
    likes: 1250,
  },
  {
    id: 2,
    name: 'Linda Wang',
    title: '首席造型师 · 8年经验',
    specialty: ['染烫专家', '护发护理', '头皮管理'],
    avatar: '/images/stylist-linda.jpg',
    rating: 4.9,
    fansCount: '1.9k',
    likes: 980,
  },
  {
    id: 3,
    name: 'Tony Zhang',
    title: '资深设计师 · 10年经验',
    specialty: ['男士精剪', '商务造型', '渐变雕刻'],
    avatar: '/images/stylist-tony.jpg',
    rating: 4.9,
    fansCount: '3.1k',
    likes: 1580,
  },
];

const stats = [
  { label: '全国门店数', value: '128+', icon: MapPin },
  { label: '注册会员数', value: '50万+', icon: Users },
  { label: '年服务人次', value: '200万+', icon: Scissors },
  { label: '客户满意度', value: '99.2%', icon: Star },
];

export default function HomePage() {
  return (
    <main className="min-h-screen bg-white">
      {/* ===== Hero区 - 品牌首屏 ===== */}
      <section className="relative h-[90vh] min-h-[600px] flex items-center justify-center overflow-hidden">
        {/* 背景视频/图片 */}
        <div className="absolute inset-0 z-0">
          <div className="absolute inset-0 bg-gradient-to-br from-[#1A1A1A]/90 via-[#1A1A1A]/70 to-[#C8A882]/30 z-10" />
          <Image
            src="https://images.unsplash.com/photo-1560066984-4d5db296fc50?w=1920&q=80"
            alt="HairCut Salon Interior"
            fill
            className="object-cover"
            priority
          />
        </div>

        {/* Hero内容 */}
        <div className="relative z-20 text-center px-6 max-w-4xl mx-auto">
          <p className="text-[#C8A882] font-medium text-lg md:text-xl mb-4 tracking-widest uppercase">
            Premium Hair Experience
          </p>
          <h1 className="text-5xl md:text-7xl lg:text-8xl font-bold text-white leading-tight mb-6">
            发现你的
            <br />
            <span className="bg-gradient-to-r from-[#C8A882] to-[#D4B896] bg-clip-text text-transparent">
              独特风格
            </span>
          </h1>
          <p className="text-lg md:text-xl text-gray-300 max-w-2xl mx-auto mb-10 leading-relaxed">
            HairCut 连锁理发品牌，汇聚顶尖造型团队，
            以匠心技艺为您打造专属形象。全国128家门店，期待您的光临。
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              href="/stores"
              className="group inline-flex items-center justify-center gap-2 bg-[#C8A882] hover:bg-[#B8956A] text-white font-semibold px-8 py-4 rounded-full text-lg transition-all duration-300 shadow-lg shadow-[#C8A882]/25 hover:shadow-xl hover:shadow-[#C8A882]/40 hover:-translate-y-0.5"
            >
              探索附近门店
              <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
            </Link>
            <Link
              href="/stylists"
              className="inline-flex items-center justify-center gap-2 border-2 border-white/30 hover:border-white/60 text-white font-semibold px-8 py-4 rounded-full text-lg transition-all duration-300 backdrop-blur-sm hover:bg-white/10"
            >
              认识明星理发师
            </Link>
          </div>

          {/* 滚动引导动画 */}
          <div className="absolute bottom-12 left-1/2 -translate-x-1/2 animate-bounce">
            <div className="w-6 h-10 rounded-full border-2 border-white/40 flex justify-center pt-2">
              <div className="w-1.5 h-3 bg-white/60 rounded-full animate-pulse" />
            </div>
          </div>
        </div>
      </section>

      {/* ===== 数据展示区 ===== */}
      <section className="py-20 bg-[#F7F8FA]">
        <div className="container mx-auto px-6">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
            {stats.map((stat, index) => (
              <div key={index} className="text-center space-y-3 group">
                <div className="inline-flex items-center justify-center w-16 h-16 bg-[#C8A882]/10 rounded-2xl group-hover:bg-[#C8A882]/20 transition-colors">
                  <stat.icon className="w-8 h-8 text-[#C8A882]" />
                </div>
                <div className="text-4xl md:text-5xl font-bold text-[#1A1A1A]">{stat.value}</div>
                <div className="text-[#6B7280] text-lg">{stat.label}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ===== 精选门店区 ===== */}
      <section className="py-24 bg-white">
        <div className="container mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-4xl md:text-5xl font-bold text-[#1A1A1A] mb-4">精选门店</h2>
            <p className="text-xl text-[#6B7280] max-w-2xl mx-auto">
              每一家门店都经过精心设计，为您提供舒适优雅的体验空间
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {featuredStores.map((store) => (
              <Link
                key={store.id}
                href={`/stores/${store.id}`}
                className="group bg-white rounded-3xl overflow-hidden shadow-card hover:shadow-card-hover transition-all duration-300 hover:-translate-y-2"
              >
                {/* 图片区 */}
                <div className="relative h-64 overflow-hidden">
                  <Image
                    src={`https://placehold.co/600x400/C8A882/white?text=Store+${store.id}`}
                    alt={store.name}
                    fill
                    className="object-cover group-hover:scale-110 transition-transform duration-500"
                  />
                  <div className="absolute top-4 left-4 flex items-center gap-2 bg-white/90 backdrop-blur-sm px-3 py-1.5 rounded-full">
                    <Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
                    <span className="font-semibold text-sm text-[#1A1A1A]">{store.rating}</span>
                    <span className="text-xs text-[#6B7280]">({store.reviews}条评价)</span>
                  </div>
                </div>

                {/* 信息区 */}
                <div className="p-6 space-y-4">
                  <h3 className="text-xl font-bold text-[#1A1A1A] group-hover:text-[#C8A882] transition-colors line-clamp-1">
                    {store.name}
                  </h3>
                  
                  <div className="flex items-center gap-2 text-[#6B7280] text-sm">
                    <MapPin className="w-4 h-4 shrink-0" />
                    <span>{store.city} · {store.address}</span>
                  </div>

                  <div className="flex items-center justify-between pt-2 border-t border-[#F3F4F6]">
                    <div className="flex gap-2 flex-wrap">
                      {store.tags.slice(0, 2).map((tag, i) => (
                        <span key={i} className="text-xs px-2 py-1 bg-[#FBF7F0] text-[#B8956A] rounded-md">
                          {tag}
                        </span>
                      ))}
                    </div>
                    <span className="font-bold text-[#C8A882]">¥{store.avgPrice}/人</span>
                  </div>
                </div>
              </Link>
            ))}
          </div>

          <div className="text-center mt-12">
            <Link
              href="/stores"
              className="inline-flex items-center gap-2 text-[#C8A882] font-semibold text-lg hover:text-[#B8956A] transition-colors group"
            >
              查看全部 128 家门店
              <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
            </Link>
          </div>
        </div>
      </section>

      {/* ===== 明星理发师区 ===== */}
      <section className="py-24 bg-[#1A1A1A] text-white relative overflow-hidden">
        {/* 背景装饰 */}
        <div className="absolute top-0 right-0 w-[600px] h-[600px] bg-[#C8A882]/5 rounded-full blur-3xl" />

        <div className="container mx-auto px-6 relative z-10">
          <div className="text-center mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-4">明星理发师</h2>
            <p className="text-xl text-gray-400 max-w-2xl mx-auto">
              汇聚行业顶尖人才，平均从业经验超过8年，累计服务超200万次
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {topStylists.map((stylist) => (
              <div
                key={stylist.id}
                className="bg-white/5 backdrop-blur-lg rounded-3xl p-8 text-center border border-white/10 hover:border-[#C8A882]/30 transition-all duration-300 hover:bg-white/10 group"
              >
                {/* 头像 */}
                <div className="relative inline-block mb-6">
                  <div className="w-32 h-32 rounded-full bg-gradient-to-br from-[#C8A882] to-[#B8956A] p-1 mx-auto">
                    <div className="w-full h-full rounded-full bg-[#2a2a2a] flex items-center justify-center text-4xl font-bold text-[#C8A882]">
                      {stylist.name[0]}
                    </div>
                  </div>
                  {/* 在线状态指示 */}
                  <div className="absolute bottom-2 right-2 w-5 h-5 bg-green-500 rounded-full border-4 border-[#1A1A1A]" />
                </div>

                <h3 className="text-2xl font-bold mb-2">{stylist.name}</h3>
                <p className="text-[#C8A882] mb-4">{stylist.title}</p>

                {/* 擅长标签 */}
                <div className="flex flex-wrap justify-center gap-2 mb-6">
                  {stylist.specialty.map((skill, i) => (
                    <span key={i} className="text-sm px-3 py-1 bg-white/10 rounded-full text-gray-300">
                      {skill}
                    </span>
                  ))}
                </div>

                {/* 统计 */}
                <div className="flex justify-center gap-8 text-sm text-gray-400 mb-6">
                  <div>
                    <div className="text-xl font-bold text-white">{stylist.fansCount}</div>
                    <div>粉丝</div>
                  </div>
                  <div>
                    <div className="text-xl font-bold text-white">{stylist.likes}</div>
                    <div>获赞</div>
                  </div>
                  <div>
                    <div className="text-xl font-bold text-yellow-400">{stylist.rating}</div>
                    <div>评分</div>
                  </div>
                </div>

                <button className="w-full py-3 bg-transparent border-2 border-[#C8A882] text-[#C8A882] rounded-full font-semibold hover:bg-[#C8A882] hover:text-white transition-all duration-300">
                  立即预约
                </button>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ===== CTA行动召唤区 ===== */}
      <section className="py-24 bg-gradient-to-br from-[#C8A882] to-[#B8956A] text-white relative overflow-hidden">
        <div className="absolute inset-0 opacity-10">
          <div className="absolute top-10 left-10 w-40 h-40 border border-white rounded-full" />
          <div className="absolute bottom-10 right-10 w-64 h-64 border border-white rounded-full" />
          <div className="absolute top-1/2 left-1/3 w-24 h-24 border border-white rounded-full" />
        </div>

        <div className="container mx-auto px-6 text-center relative z-10">
          <Sparkles className="w-12 h-12 mx-auto mb-6 opacity-80" />
          <h2 className="text-4xl md:text-5xl font-bold mb-6">准备好改变了吗？</h2>
          <p className="text-xl text-white/80 max-w-2xl mx-auto mb-10">
            立即预约，开启您的焕新之旅。新用户首次预约享8折优惠！
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              href="/stores"
              className="inline-flex items-center justify-center gap-2 bg-white text-[#C8A882] font-bold px-10 py-4 rounded-full text-lg hover:bg-gray-100 transition-colors shadow-lg"
            >
              预约到店服务
            </Link>
            <Link
              href="/franchise"
              className="inline-flex items-center justify-center gap-2 bg-transparent border-2 border-white text-white font-bold px-10 py-4 rounded-full text-lg hover:bg-white/10 transition-colors"
            >
              了解加盟合作
            </Link>
          </div>
        </div>
      </section>
    </main>
  );
}
