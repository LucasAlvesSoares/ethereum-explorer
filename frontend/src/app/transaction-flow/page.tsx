'use client';

import { useState, useEffect, useRef } from 'react';
import * as d3 from 'd3';
import { TransactionNode, TransactionLink, TransactionFlowData } from '@/types';
import { isValidEthereumAddress } from '@/utils/validation';

export default function TransactionFlowPage() {
  const [flowData, setFlowData] = useState<TransactionFlowData>({ nodes: [], links: [] });
  const [loading, setLoading] = useState(false);
  const [selectedAddress, setSelectedAddress] = useState('');
  const [searchAddress, setSearchAddress] = useState('');
  const [error, setError] = useState<string | null>(null);
  const svgRef = useRef<SVGSVGElement>(null);


  const fetchTransactionFlow = async (address: string) => {
    if (!address) return;
    
    // Clear previous error
    setError(null);
    
    // Validate address format on client side first
    if (!isValidEthereumAddress(address)) {
      setError('Invalid Ethereum address format. Address must be 42 characters long and start with 0x followed by hexadecimal characters.');
      setFlowData({ nodes: [], links: [] });
      setSelectedAddress('');
      return;
    }
    
    setLoading(true);
    try {
      const response = await fetch(`/api/v1/transaction-flow/${address}`);
      if (response.ok) {
        const data: TransactionFlowData = await response.json();
        setFlowData(data);
        setSelectedAddress(address);
        setError(null);
      } else if (response.status === 404) {
        // Handle addresses with no transaction data
        const errorData = await response.json().catch(() => ({ error: 'No transaction data found' }));
        setError(`No transaction data found for this address. This address may be new or have no transaction history.`);
        setFlowData({ nodes: [], links: [] });
        setSelectedAddress('');
      } else {
        // Handle other API errors
        const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
        setError(errorData.error || `Server error: ${response.status}`);
        setFlowData({ nodes: [], links: [] });
        setSelectedAddress('');
      }
    } catch (error) {
      console.error('Error fetching transaction flow:', error);
      setError('Network error: Unable to connect to the server. Please try again.');
      setFlowData({ nodes: [], links: [] });
      setSelectedAddress('');
    } finally {
      setLoading(false);
    }
  };


  useEffect(() => {
    if (!flowData.nodes.length || !svgRef.current) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll('*').remove();

    const width = 800;
    const height = 600;
    const margin = { top: 20, right: 20, bottom: 20, left: 20 };

    svg.attr('width', width).attr('height', height);

    const g = svg.append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    // Create scales
    const nodeScale = d3.scaleLinear()
      .domain(d3.extent(flowData.nodes, d => d.value) as [number, number])
      .range([5, 20]);

    const linkScale = d3.scaleLinear()
      .domain(d3.extent(flowData.links, d => d.value) as [number, number])
      .range([1, 8]);

    // Create force simulation
    const simulation = d3.forceSimulation(flowData.nodes as any)
      .force('link', d3.forceLink(flowData.links).id((d: any) => d.id).distance(100))
      .force('charge', d3.forceManyBody().strength(-300))
      .force('center', d3.forceCenter((width - margin.left - margin.right) / 2, (height - margin.top - margin.bottom) / 2))
      .force('collision', d3.forceCollide().radius((d: any) => nodeScale(d.value) + 2));

    // Create arrow markers
    svg.append('defs').selectAll('marker')
      .data(['arrow'])
      .enter().append('marker')
      .attr('id', 'arrow')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 15)
      .attr('refY', 0)
      .attr('markerWidth', 6)
      .attr('markerHeight', 6)
      .attr('orient', 'auto')
      .append('path')
      .attr('d', 'M0,-5L10,0L0,5')
      .attr('fill', '#666');

    // Create links
    const link = g.append('g')
      .selectAll('line')
      .data(flowData.links)
      .enter().append('line')
      .attr('stroke', '#999')
      .attr('stroke-opacity', 0.6)
      .attr('stroke-width', (d: any) => linkScale(d.value))
      .attr('marker-end', 'url(#arrow)');

    // Create nodes
    const node = g.append('g')
      .selectAll('circle')
      .data(flowData.nodes)
      .enter().append('circle')
      .attr('r', (d: any) => nodeScale(d.value))
      .attr('fill', (d: any) => d.type === 'contract' ? '#ff6b6b' : '#4ecdc4')
      .attr('stroke', '#fff')
      .attr('stroke-width', 2)
      .style('cursor', 'pointer')
      .call(d3.drag<any, any>()
        .on('start', (event, d: any) => {
          if (!event.active) simulation.alphaTarget(0.3).restart();
          d.fx = d.x;
          d.fy = d.y;
        })
        .on('drag', (event, d: any) => {
          d.fx = event.x;
          d.fy = event.y;
        })
        .on('end', (event, d: any) => {
          if (!event.active) simulation.alphaTarget(0);
          d.fx = null;
          d.fy = null;
        }));

    // Add labels
    const labels = g.append('g')
      .selectAll('text')
      .data(flowData.nodes)
      .enter().append('text')
      .text((d: any) => d.label || `${d.address.slice(0, 6)}...${d.address.slice(-4)}`)
      .attr('font-size', 10)
      .attr('font-family', 'monospace')
      .attr('text-anchor', 'middle')
      .attr('dy', -25)
      .attr('fill', '#333');

    // Add tooltips
    const tooltip = d3.select('body').append('div')
      .attr('class', 'tooltip')
      .style('opacity', 0)
      .style('position', 'absolute')
      .style('background', 'rgba(0, 0, 0, 0.8)')
      .style('color', 'white')
      .style('padding', '8px')
      .style('border-radius', '4px')
      .style('font-size', '12px')
      .style('pointer-events', 'none');

    node
      .on('mouseover', (event, d: any) => {
        tooltip.transition().duration(200).style('opacity', .9);
        tooltip.html(`
          <strong>Address:</strong> ${d.address}<br/>
          <strong>Type:</strong> ${d.type}<br/>
          <strong>Value:</strong> ${d.value.toFixed(2)} ETH
        `)
          .style('left', (event.pageX + 10) + 'px')
          .style('top', (event.pageY - 28) + 'px');
      })
      .on('mouseout', () => {
        tooltip.transition().duration(500).style('opacity', 0);
      });

    // Update positions on simulation tick
    simulation.on('tick', () => {
      link
        .attr('x1', (d: any) => d.source.x)
        .attr('y1', (d: any) => d.source.y)
        .attr('x2', (d: any) => d.target.x)
        .attr('y2', (d: any) => d.target.y);

      node
        .attr('cx', (d: any) => d.x)
        .attr('cy', (d: any) => d.y);

      labels
        .attr('x', (d: any) => d.x)
        .attr('y', (d: any) => d.y);
    });

    return () => {
      tooltip.remove();
    };
  }, [flowData]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchAddress.trim()) {
      fetchTransactionFlow(searchAddress.trim());
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Transaction Flow Visualization</h1>
        <p className="text-gray-600 mb-6">
          Explore the flow of transactions between addresses with interactive network graphs.
        </p>

        {/* Search Form */}
        <form onSubmit={handleSearch} className="mb-6">
          <div className="flex gap-4">
            <input
              type="text"
              value={searchAddress}
              onChange={(e) => setSearchAddress(e.target.value)}
              placeholder="Enter Ethereum address (0x...)"
              className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <button
              type="submit"
              disabled={loading}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? 'Loading...' : 'Visualize'}
            </button>
          </div>
        </form>

        {/* Quick Examples */}
        <div className="mb-6">
          <p className="text-sm text-gray-600 mb-2">Try these example addresses:</p>
          <div className="flex flex-wrap gap-2">
            {[
              '0xdAC17F958D2ee523a2206206994597C13D831ec7', // USDT Contract
              '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984', // Uniswap Token
              '0x28C6c06298d514Db089934071355E5743bf21d60'  // Binance Hot Wallet
            ].map((addr) => (
              <button
                key={addr}
                onClick={() => {
                  setSearchAddress(addr);
                  fetchTransactionFlow(addr);
                }}
                className="px-3 py-1 text-xs bg-gray-100 hover:bg-gray-200 rounded-full text-gray-700"
              >
                {addr.slice(0, 6)}...{addr.slice(-4)}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Error Display */}
      {error && (
        <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800">Validation Error</h3>
              <div className="mt-2 text-sm text-red-700">
                {error}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Visualization */}
      {loading ? (
        <div className="flex justify-center items-center h-96">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      ) : flowData.nodes.length > 0 ? (
        <div className="bg-white rounded-lg shadow-lg p-6">
          <div className="mb-4">
            <h2 className="text-xl font-semibold text-gray-900 mb-2">
              Transaction Flow for {selectedAddress.slice(0, 6)}...{selectedAddress.slice(-4)}
            </h2>
            <div className="flex items-center gap-6 text-sm text-gray-600">
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-teal-400"></div>
                <span>Regular Address</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-red-400"></div>
                <span>Smart Contract</span>
              </div>
              <div className="text-xs">
                Drag nodes to rearrange • Hover for details
              </div>
            </div>
          </div>
          
          <div className="border rounded-lg overflow-hidden">
            <svg ref={svgRef} className="w-full"></svg>
          </div>

          {/* Statistics */}
          <div className="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="bg-gray-50 p-4 rounded-lg">
              <div className="text-2xl font-bold text-gray-900">{flowData.nodes.length}</div>
              <div className="text-sm text-gray-600">Connected Addresses</div>
            </div>
            <div className="bg-gray-50 p-4 rounded-lg">
              <div className="text-2xl font-bold text-gray-900">{flowData.links.length}</div>
              <div className="text-sm text-gray-600">Transactions</div>
            </div>
            <div className="bg-gray-50 p-4 rounded-lg">
              <div className="text-2xl font-bold text-gray-900">
                {flowData.links.reduce((sum, link) => sum + link.value, 0).toFixed(2)}
              </div>
              <div className="text-sm text-gray-600">Total ETH Flow</div>
            </div>
          </div>
        </div>
      ) : (
        <div className="text-center py-12">
          <div className="text-gray-500 mb-4">
            <svg className="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No Transaction Flow Data</h3>
          <p className="text-gray-600">Enter an Ethereum address above to visualize its transaction flow.</p>
        </div>
      )}
    </div>
  );
}
