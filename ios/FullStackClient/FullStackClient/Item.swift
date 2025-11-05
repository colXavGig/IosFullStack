//
//  Item.swift
//  FullStackClient
//
//  Created by Xavier Giguerre on 2025-11-02.
//

import Foundation
import SwiftData

@Model
final class Item {
    var timestamp: Date
    
    init(timestamp: Date) {
        self.timestamp = timestamp
    }
}
