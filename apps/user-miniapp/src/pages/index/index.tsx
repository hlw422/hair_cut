import { Component } from 'react';
import { View, Swiper, SwiperItem, Image, Text } from '@tarojs/components';
import { navigateTo } from '@tarojs/taro';
import './index.module.less';

// Mock数据（实际应从API获取）
const mockBanners = [
  { id: 1, image: 'https://placehold.co/750x320/C8A882/white?text=New+Year+Sale+2024', link: '' },
  { id: 2, image: 'https://placehold.co/750x320/1A1A1A/C8A882?text=VIP+Member+Exclusive', link: '' },
  { id: 3, image: 'https://placehold.co/750x320/F7F8FA/C8A882?text=Ai+Hairstyle+Analysis', link: '/pages/ai-hairstyle/ai-hairstyle' },
];

const quickEntries = [
  { icon: '📍', label: '附近门店', path: '/pages/store-list/store-list' },
  { icon: '⭐', label: '明星理发师', path: '/pages/store-detail/stylist-list' },
  { icon: '🎁', label: '新人礼包', path: '/pages/coupon/new-user-coupon' },
  { icon: '🎟️', label: '优惠券', path: '/pages/coupon/coupon-list' },
  { icon: '💰', label: '积分商城', path: '/pages/points-mall/points-mall' },
  { icon: '✨', label: 'AI测发型', path: '/pages/ai-hairstyle/ai-hairstyle' },
  { icon: '📋', label: '我的订单', path: '/pages/order/order-list' },
  { icon: '💬', label: '消息中心', path: '/pages/messages/messages' },
];

const recommendedStores = [
  {
    id: 1,
    name: 'HairCut 精品沙龙（静安店）',
    logo: 'https://placehold.co/200x200/C8A882/white?text=HC',
    distance: '850m',
    rating: 4.8,
    reviews: 256,
    avgPrice: 168,
    tags: ['环境优雅', '明星理发师'],
    isOpen: true,
  },
  {
    id: 2,
    name: 'HairCut 造型中心（徐汇店）',
    logo: 'https://placehold.co/200x200/C8A882/white?text=HC',
    distance: '2.3km',
    rating: 4.9,
    reviews: 389,
    avgPrice: 198,
    tags: ['人气TOP1', '免费停车'],
    isOpen: true,
  },
  {
    id: 3,
    name: 'HairCut 概念店（浦东）',
    logo: 'https://placehold.co/200x200/C8A882/white?text=HC',
    distance: '5.6km',
    rating: 4.7,
    reviews: 178,
    avgPrice: 158,
    tags: ['新店开业', '8折优惠'],
    isOpen: true,
  },
];

const hotStylists = [
  { id: 1, name: 'Kevin老师', avatar: '', title: '艺术总监', exp: 12, rating: 4.9, fans: 2890 },
  { id: 2, name: 'Linda老师', avatar: '', title: '首席造型师', exp: 8, rating: 4.8, fans: 1920 },
  { id: 3, name: 'Tony老师', avatar: '', title: '资深设计师', exp: 10, rating: 4.9, fans: 3150 },
];

const hotPackages = [
  { id: 1, name: '新年焕新套餐', cover: 'https://placehold.co/375x280/C8A882/white?text=Package+1', price: 298, originalPrice: 468, sold: 520 },
  { id: 2, name: '染烫护理三合一', cover: 'https://placehold.co/375x280/B8956A/white?text=Package+2', price: 568, originalPrice: 888, sold: 380 },
  { id: 3, name: '男士精剪+修面', cover: 'https://placehold.co/375x280/1A1A1A/C8A882?text=Package+3', price: 128, originalPrice: 188, sold: 892 },
  { id: 4, name: '儿童专属造型', cover: 'https://placehold.co/375x280/F7F8FA/C8A882?text=Package+4', price: 98, originalPrice: 148, sold: 456 },
];

interface IndexState {
  city: string;
}

class Index extends Component<{}, IndexState> {
  constructor(props) {
    super(props);
    this.state = {
      city: '上海市',
    };
  }

  componentDidMount() {
    // TODO: 获取用户位置并设置城市
    // Taro.getLocation().then(res => this.getCityFromLocation(res));
  }

  handleCityChange = () => {
    // TODO: 打开城市选择器
  };

  navigateTo = (url: string) => {
    navigateTo({ url });
  };

  render() {
    const { city } = this.state;

    return (
      <View className="index-page">
        {/* 顶部城市定位栏 */}
        <View className="header-bar">
          <View className="city-selector" onClick={this.handleCityChange}>
            <Text className="city-icon">📍</Text>
            <Text className="city-name">{city}</Text>
            <Text className="arrow">▼</Text>
          </View>
          <View className="search-bar" onClick={() => this.navigateTo('/pages/store-list/store-list')}>
            <Text className="search-placeholder">搜索门店、理发师、服务...</Text>
          </View>
        </View>

        {/* Banner轮播图 */}
        <Swiper
          className="banner-swiper"
          indicatorColor="rgba(255,255,255,0.5)"
          indicatorActiveColor="#C8A882"
          circular
          autoplay
          interval={4000}
        >
          {mockBanners.map((banner) => (
            <SwiperItem key={banner.id}>
              <Image
                className="banner-image"
                src={banner.image}
                mode="aspectFill"
                onClick={() => banner.link && this.navigateTo(banner.link)}
              />
            </SwiperItem>
          ))}
        </Swiper>

        {/* 快捷入口 - 8宫格 */}
        <View className="quick-entries">
          {quickEntries.map((entry, index) => (
            <View
              key={index}
              className="entry-item"
              onClick={() => this.navigateTo(entry.path)}
            >
              <View className="entry-icon">{entry.icon}</View>
              <Text className="entry-label">{entry.label}</Text>
            </View>
          ))}
        </View>

        {/* 今日推荐 - 附近门店横向滑动 */}
        <View className="section recommended-stores">
          <View className="section-header">
            <Text className="section-title">🔥 今日推荐</Text>
            <Text className="more-link" onClick={() => this.navigateTo('/pages/store-list/store-list')}>查看更多 ›</Text>
          </View>
          <ScrollView scrollX className="store-scroll">
            <View className="store-list-horizontal">
              {recommendedStores.map((store) => (
                <View
                  key={store.id}
                  className="store-card"
                  onClick={() => this.navigateTo(`/pages/store-detail/store-detail?id=${store.id}`)}
                >
                  <Image className="store-logo" src={store.logo} mode="aspectFill" />
                  <Text className="store-name">{store.name.split('（')[0]}</Text>
                  <View className="store-info">
                    <Text className="distance">{store.distance}</Text>
                    <Text className="rating">★ {store.rating}</Text>
                  </View>
                  <View className="tags">
                    {store.tags.map((tag, i) => (
                      <Text key={i} className="tag">{tag}</Text>
                    ))}
                  </View>
                  <View className={`status-badge ${store.isOpen ? 'open' : 'closed'}`}>
                    {store.isOpen ? '营业中' : '休息中'}
                  </View>
                </View>
              ))}
            </View>
          </ScrollView>
        </View>

        {/* 明星理发师 */}
        <View className="section stylists-section">
          <View className="section-header">
            <Text className="section-title">⭐ 明星理发师</Text>
            <Text className="more-link">更多 ›</Text>
          </View>
          <ScrollView scrollX className="stylist-scroll">
            <View className="stylist-list-horizontal">
              {hotStylists.map((stylist) => (
                <View key={stylist.id} className="stylist-card">
                  <View className="avatar-circle">
                    <Text className="avatar-text">{stylist.name[0]}</Text>
                  </View>
                  <Text className="stylist-name">{stylist.name}</Text>
                  <Text className="stylist-title">{stylist.title}</Text>
                  <View className="stylist-stats">
                    <Text>从业{stylist.exp}年 · ★{stylist.rating}</Text>
                  </View>
                  <Text className="fan-count">{(stylist.fans / 1000).toFixed(1)}k粉丝</Text>
                  <View className="book-btn" onClick={() => this.navigateTo('/pages/appointment/select-store')}>
                    预约
                  </View>
                </View>
              ))}
            </View>
          </ScrollView>
        </View>

        {/* 热门套餐 - 瀑布流双列 */}
        <View className="section packages-section">
          <View className="section-header">
            <Text className="section-title">🎁 热门套餐</Text>
            <Text className="more-link">全部 ›</Text>
          </View>
          <View className="packages-grid">
            {hotPackages.map((pkg) => (
              <View key={pkg.id} className="package-card">
                <Image className="package-cover" src={pkg.cover} mode="aspectFill" />
                <View className="package-info">
                  <Text className="package-name">{pkg.name}</Text>
                  <View className="package-price-row">
                    <Text className="price-symbol">¥</Text>
                    <Text className="price-value">{pkg.price}</Text>
                    <Text className="original-price">¥{pkg.originalPrice}</Text>
                  </View>
                  <Text className="sold-count">已售 {pkg.sold}</Text>
                </View>
              </View>
            ))}
          </View>
        </View>
      </View>
    );
  }
}

export default Index;
